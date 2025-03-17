package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	db "texApi/database"
	"texApi/internal/dto"
	"texApi/pkg/utils"
)

// RecordGPSLocation records a single GPS location
func RecordGPSLocation(ctx *gin.Context) {
	var location dto.GPSLocationCreate
	companyID := ctx.MustGet("companyID").(int)
	role := ctx.MustGet("role").(string)

	if err := ctx.ShouldBindJSON(&location); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	// Verify the vehicle and driver belong to the company
	if role != "admin" {
		var count int
		err := db.DB.QueryRow(
			context.Background(),
			`SELECT COUNT(*) FROM tbl_vehicle v
			JOIN tbl_driver d ON d.id = $1 AND d.company_id = $3
			WHERE v.id = $2 AND v.company_id = $3`,
			location.DriverID, location.VehicleID, companyID,
		).Scan(&count)
		if err != nil || count == 0 {
			ctx.JSON(http.StatusForbidden, utils.FormatErrorResponse("Permission denied", "Vehicle or driver does not belong to your company"))
			return
		}
	}

	// Insert the device if it doesn't exist
	deviceQuery := `
		INSERT INTO tbl_gps_device (vehicle_id, driver_id, device_id, device_type, last_ping) 
		VALUES ($1, $2, $3, 'api', CURRENT_TIMESTAMP)
		ON CONFLICT (device_id) DO UPDATE SET 
		vehicle_id = $1, driver_id = $2, last_ping = CURRENT_TIMESTAMP
	`
	_, err := db.DB.Exec(
		context.Background(),
		deviceQuery,
		location.VehicleID, location.DriverID, location.DeviceID,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error registering device", err.Error()))
		return
	}

	// Insert the GPS location
	query := `
		INSERT INTO tbl_gps_location (
			vehicle_id, driver_id, latitude, longitude, 
			altitude, speed, direction, accuracy, 
			location_time, meta
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		) RETURNING id, uuid
	`

	var altitude, speed, direction, accuracy interface{} = nil, nil, nil, nil
	if location.Altitude != nil {
		altitude = *location.Altitude
	}
	if location.Speed != nil {
		speed = *location.Speed
	}
	if location.Direction != nil {
		direction = *location.Direction
	}
	if location.Accuracy != nil {
		accuracy = *location.Accuracy
	}

	// If meta is empty, use empty JSON object
	meta := location.Meta
	if len(meta) == 0 {
		meta = []byte("{}")
	}

	var locationID int
	var locationUUID string
	err = db.DB.QueryRow(
		context.Background(),
		query,
		location.VehicleID, location.DriverID,
		location.Latitude, location.Longitude,
		altitude, speed, direction, accuracy,
		location.LocationTime, meta,
	).Scan(&locationID, &locationUUID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error recording location", err.Error()))
		return
	}

	//// Check for geofence events asynchronously
	//go checkGeofenceEvents(location)

	ctx.JSON(http.StatusCreated, utils.FormatResponse("Location recorded successfully", gin.H{
		"location_id": locationID,
		"uuid":        locationUUID,
	}))
}

// RecordBatchGPSLocation records multiple GPS locations in a batch
func RecordBatchGPSLocation(ctx *gin.Context) {
	var batchLocation dto.GPSBatchLocationCreate
	companyID := ctx.MustGet("companyID").(int)
	role := ctx.MustGet("role").(string)

	if err := ctx.ShouldBindJSON(&batchLocation); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	// Begin transaction
	tx, err := db.DB.Begin(context.Background())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Database error", err.Error()))
		return
	}
	defer tx.Rollback(context.Background())

	// Verify permissions for each vehicle/driver pair
	if role != "admin" {
		for _, location := range batchLocation.Locations {
			var count int
			err := tx.QueryRow(
				context.Background(),
				`SELECT COUNT(*) FROM tbl_vehicle v
				JOIN tbl_driver d ON d.id = $1 AND d.company_id = $3
				WHERE v.id = $2 AND v.company_id = $3`,
				location.DriverID, location.VehicleID, companyID,
			).Scan(&count)
			if err != nil || count == 0 {
				ctx.JSON(http.StatusForbidden, utils.FormatErrorResponse("Permission denied", "Vehicle or driver does not belong to your company"))
				return
			}
		}
	}

	// Update device info
	if len(batchLocation.Locations) > 0 {
		latestLocation := batchLocation.Locations[len(batchLocation.Locations)-1]
		deviceQuery := `
			INSERT INTO tbl_gps_device (vehicle_id, driver_id, device_id, device_type, last_ping) 
			VALUES ($1, $2, $3, 'api', CURRENT_TIMESTAMP)
			ON CONFLICT (device_id) DO UPDATE SET 
			vehicle_id = $1, driver_id = $2, last_ping = CURRENT_TIMESTAMP
		`
		_, err := tx.Exec(
			context.Background(),
			deviceQuery,
			latestLocation.VehicleID, latestLocation.DriverID, batchLocation.DeviceID,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error registering device", err.Error()))
			return
		}
	}

	// Insert all locations
	locationIDs := make([]int, 0, len(batchLocation.Locations))
	for _, location := range batchLocation.Locations {
		query := `
			INSERT INTO tbl_gps_location (
				vehicle_id, driver_id, latitude, longitude, 
				altitude, speed, direction, accuracy, 
				location_time, meta
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
			) RETURNING id
		`

		var altitude, speed, direction, accuracy interface{} = nil, nil, nil, nil
		if location.Altitude != nil {
			altitude = *location.Altitude
		}
		if location.Speed != nil {
			speed = *location.Speed
		}
		if location.Direction != nil {
			direction = *location.Direction
		}
		if location.Accuracy != nil {
			accuracy = *location.Accuracy
		}

		// If meta is empty, use empty JSON object
		meta := location.Meta
		if len(meta) == 0 {
			meta = []byte("{}")
		}

		var locationID int
		err = tx.QueryRow(
			context.Background(),
			query,
			location.VehicleID, location.DriverID,
			location.Latitude, location.Longitude,
			altitude, speed, direction, accuracy,
			location.LocationTime, meta,
		).Scan(&locationID)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error recording location batch", err.Error()))
			return
		}

		locationIDs = append(locationIDs, locationID)
	}

	// Commit transaction
	if err := tx.Commit(context.Background()); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error committing transaction", err.Error()))
		return
	}

	//// Check for geofence events asynchronously
	//for _, location := range batchLocation.Locations {
	//	go checkGeofenceEvents(location)
	//}

	ctx.JSON(http.StatusCreated, utils.FormatResponse("Batch locations recorded successfully", gin.H{
		"count":     len(locationIDs),
		"first_id":  locationIDs[0],
		"device_id": batchLocation.DeviceID,
	}))
}

// GetCurrentVehicleLocation retrieves the most recent location for a vehicle
func GetCurrentVehicleLocation(ctx *gin.Context) {
	vehicleID, err := strconv.Atoi(ctx.Param("vehicle_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid vehicle ID", err.Error()))
		return
	}

	companyID := ctx.MustGet("companyID").(int)
	role := ctx.MustGet("role").(string)

	// Verify the vehicle belongs to the company
	if role != "admin" {
		var count int
		err := db.DB.QueryRow(
			context.Background(),
			`SELECT COUNT(*) FROM tbl_vehicle WHERE id = $1 AND company_id = $2 AND active = 1 AND deleted = 0`,
			vehicleID, companyID,
		).Scan(&count)
		if err != nil || count == 0 {
			ctx.JSON(http.StatusForbidden, utils.FormatErrorResponse("Permission denied", "Vehicle does not belong to your company"))
			return
		}
	}

	query := `
		SELECT 
			gl.id, gl.uuid, gl.vehicle_id, gl.driver_id, 
			gl.latitude, gl.longitude, gl.altitude, gl.speed, 
			gl.direction, gl.accuracy, gl.location_time, gl.created_at, gl.meta,
			d.first_name || ' ' || d.last_name AS driver_name,
			v.numberplate AS vehicle_plate
		FROM tbl_gps_location gl
		JOIN tbl_driver d ON gl.driver_id = d.id
		JOIN tbl_vehicle v ON gl.vehicle_id = v.id
		WHERE gl.vehicle_id = $1
		ORDER BY gl.location_time DESC
		LIMIT 1
	`

	var location dto.GPSLocationDetails
	err = pgxscan.Get(context.Background(), db.DB, &location, query, vehicleID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("No location data found for this vehicle", err.Error()))
		return
	}

	//TODO:
	//// Calculate the location age
	//locationAge := time.Since(location.LocationTime)
	//location.LocationAge = formatDuration(locationAge)

	ctx.JSON(http.StatusOK, utils.FormatResponse("Current vehicle location", location))
}

// GetVehicleLocationHistory retrieves historical location data for a vehicle
func GetVehicleLocationHistory(ctx *gin.Context) {
	vehicleIDStr := ctx.Query("vehicle_id")
	startTimeStr := ctx.Query("start_time")
	endTimeStr := ctx.Query("end_time")
	intervalStr := ctx.DefaultQuery("interval", "0")
	includeMetaStr := ctx.DefaultQuery("include_meta", "false")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "100"))
	offset := (page - 1) * perPage

	companyID := ctx.MustGet("companyID").(int)
	role := ctx.MustGet("role").(string)

	// Validate required parameters
	if vehicleIDStr == "" || startTimeStr == "" || endTimeStr == "" {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Missing required parameters", "vehicle_id, start_time, and end_time are required"))
		return
	}

	vehicleID, err := strconv.Atoi(vehicleIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid vehicle ID", err.Error()))
		return
	}

	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid start_time format", "Use ISO 8601 format (e.g., 2025-03-17T12:00:00Z)"))
		return
	}

	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid end_time format", "Use ISO 8601 format (e.g., 2025-03-17T12:00:00Z)"))
		return
	}

	interval, err := strconv.Atoi(intervalStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid interval", "Interval must be an integer"))
		return
	}

	includeMeta := includeMetaStr == "true"

	// Verify the vehicle belongs to the company
	if role != "admin" {
		var count int
		err := db.DB.QueryRow(
			context.Background(),
			`SELECT COUNT(*) FROM tbl_vehicle WHERE id = $1 AND company_id = $2 AND active = 1 AND deleted = 0`,
			vehicleID, companyID,
		).Scan(&count)
		if err != nil || count == 0 {
			ctx.JSON(http.StatusForbidden, utils.FormatErrorResponse("Permission denied", "Vehicle does not belong to your company"))
			return
		}
	}

	// Get driver info
	var driverName string
	var driverID int
	err = db.DB.QueryRow(
		context.Background(),
		`SELECT d.id, d.first_name || ' ' || d.last_name FROM tbl_driver d
		JOIN tbl_vehicle v ON v.id = $1
		WHERE d.company_id = v.company_id AND d.active = 1 AND d.deleted = 0
		LIMIT 1`,
		vehicleID,
	).Scan(&driverID, &driverName)
	if err != nil {
		driverName = "Unknown"
		driverID = 0
	}

	// Query for total count
	var totalCount int
	countQuery := `
		SELECT COUNT(*) FROM tbl_gps_location
		WHERE vehicle_id = $1 AND location_time BETWEEN $2 AND $3
	`
	err = db.DB.QueryRow(context.Background(), countQuery, vehicleID, startTime, endTime).Scan(&totalCount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error counting locations", err.Error()))
		return
	}

	// Adjust query based on interval
	var query string
	var args []interface{}
	if interval > 0 {
		// Use time_bucket for interval sampling if using an interval
		query = `
        WITH sampled AS (
            SELECT 
                id, 
                vehicle_id, 
                driver_id, 
                latitude, 
                longitude, 
                altitude, 
                speed, 
                direction, 
                accuracy, 
                location_time,
                meta,
                ROW_NUMBER() OVER (
                    PARTITION BY (
                        EXTRACT(EPOCH FROM location_time)::integer / $4
                    ) 
                    ORDER BY location_time
                ) AS rn
            FROM tbl_gps_location
            WHERE vehicle_id = $1 
            AND location_time BETWEEN $2 AND $3
        )
        SELECT id, vehicle_id, driver_id, latitude, longitude, 
               altitude, speed, direction, accuracy, location_time, 
               CASE WHEN $5 = true THEN meta ELSE '{}'::jsonb END as meta
        FROM sampled
        WHERE rn = 1
        ORDER BY location_time
        LIMIT $6 OFFSET $7
    `
		args = []interface{}{vehicleID, startTime, endTime, interval, includeMeta, perPage, offset}
	} else {
		// Get all points with pagination
		query = `
        SELECT id, vehicle_id, driver_id, latitude, longitude, 
               altitude, speed, direction, accuracy, location_time,
               CASE WHEN $4 = true THEN meta ELSE '{}'::jsonb END as meta
        FROM tbl_gps_location
        WHERE vehicle_id = $1 
        AND location_time BETWEEN $2 AND $3
        ORDER BY location_time
        LIMIT $5 OFFSET $6
    `
		args = []interface{}{vehicleID, startTime, endTime, includeMeta, perPage, offset}
	}

	// Execute query
	var locations []dto.GPSLocationDetails
	err = pgxscan.Select(context.Background(), db.DB, &locations, query, args...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error retrieving location history", err.Error()))
		return
	}

	// Calculate analytics
	var totalDistance float64
	var avgSpeed, maxSpeed float64

	// TODO:
	//if len(locations) > 0 {
	//	// Calculate distance and speeds only if we have points
	//	distanceQuery := `
	//    SELECT
	//        COALESCE(SUM(
	//            ST_DistanceSphere(
	//                ST_MakePoint(longitude, latitude),
	//                LAG(ST_MakePoint(longitude, latitude)) OVER (ORDER BY location_time)
	//            )
	//        ) / 1000, 0) as total_distance,
	//        COALESCE(AVG(speed), 0) as avg_speed,
	//        COALESCE(MAX(speed), 0) as max_speed
	//    FROM tbl_gps_location
	//    WHERE vehicle_id = $1
	//    AND location_time BETWEEN $2 AND $3
	//    AND speed IS NOT NULL
	//`
	//	err = db.DB.QueryRow(
	//		context.Background(),
	//		distanceQuery,
	//		vehicleID, startTime, endTime,
	//	).Scan(&totalDistance, &avgSpeed, &maxSpeed)
	//
	//	if err != nil {
	//		// If error in analytics calculation, continue with zeros
	//		totalDistance, avgSpeed, maxSpeed = 0, 0, 0
	//	}
	//}

	// Create response
	response := gin.H{
		"vehicle_id":   vehicleID,
		"driver_id":    driverID,
		"driver_name":  driverName,
		"total_points": totalCount,
		"points":       locations,
		"distance":     totalDistance,
		"avg_speed":    avgSpeed,
		"max_speed":    maxSpeed,
		"pagination": gin.H{
			"page":     page,
			"per_page": perPage,
			"total":    totalCount,
			"pages":    (totalCount + perPage - 1) / perPage,
		},
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Vehicle location history", response))
}

// GetDriverLocationHistory retrieves historical location data for a driver
func GetDriverLocationHistory(ctx *gin.Context) {
	driverIDStr := ctx.Query("driver_id")
	startTimeStr := ctx.Query("start_time")
	endTimeStr := ctx.Query("end_time")
	intervalStr := ctx.DefaultQuery("interval", "0")
	includeMetaStr := ctx.DefaultQuery("include_meta", "false")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "100"))
	offset := (page - 1) * perPage

	companyID := ctx.MustGet("companyID").(int)
	role := ctx.MustGet("role").(string)

	// Validate required parameters
	if driverIDStr == "" || startTimeStr == "" || endTimeStr == "" {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Missing required parameters", "driver_id, start_time, and end_time are required"))
		return
	}

	driverID, err := strconv.Atoi(driverIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid driver ID", err.Error()))
		return
	}

	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid start_time format", "Use ISO 8601 format (e.g., 2025-03-17T12:00:00Z)"))
		return
	}

	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid end_time format", "Use ISO 8601 format (e.g., 2025-03-17T12:00:00Z)"))
		return
	}

	interval, err := strconv.Atoi(intervalStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid interval", "Interval must be an integer"))
		return
	}

	includeMeta := includeMetaStr == "true"

	// Verify the driver belongs to the company
	if role != "admin" {
		var count int
		err := db.DB.QueryRow(
			context.Background(),
			`SELECT COUNT(*) FROM tbl_driver WHERE id = $1 AND company_id = $2 AND active = 1 AND deleted = 0`,
			driverID, companyID,
		).Scan(&count)
		if err != nil || count == 0 {
			ctx.JSON(http.StatusForbidden, utils.FormatErrorResponse("Permission denied", "Driver does not belong to your company"))
			return
		}
	}

	// Get driver name
	var driverName string
	err = db.DB.QueryRow(
		context.Background(),
		`SELECT first_name || ' ' || last_name FROM tbl_driver WHERE id = $1`,
		driverID,
	).Scan(&driverName)
	if err != nil {
		driverName = "Unknown"
	}

	// Query for total count
	var totalCount int
	countQuery := `
        SELECT COUNT(*) FROM tbl_gps_location
        WHERE driver_id = $1 AND location_time BETWEEN $2 AND $3
    `
	err = db.DB.QueryRow(context.Background(), countQuery, driverID, startTime, endTime).Scan(&totalCount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error counting locations", err.Error()))
		return
	}

	// Adjust query based on interval
	var query string
	var args []interface{}
	if interval > 0 {
		// Use time_bucket for interval sampling if using an interval
		query = `
            WITH sampled AS (
                SELECT 
                    id, 
                    vehicle_id, 
                    driver_id, 
                    latitude, 
                    longitude, 
                    altitude, 
                    speed, 
                    direction, 
                    accuracy, 
                    location_time,
                    meta,
                    ROW_NUMBER() OVER (
                        PARTITION BY (
                            EXTRACT(EPOCH FROM location_time)::integer / $4
                        ) 
                        ORDER BY location_time
                    ) AS rn
                FROM tbl_gps_location
                WHERE driver_id = $1 
                AND location_time BETWEEN $2 AND $3
            )
            SELECT id, vehicle_id, driver_id, latitude, longitude, 
                   altitude, speed, direction, accuracy, location_time,
                   CASE WHEN $5 = true THEN meta ELSE '{}'::jsonb END as meta
            FROM sampled
            WHERE rn = 1
            ORDER BY location_time
            LIMIT $6 OFFSET $7
        `
		args = []interface{}{driverID, startTime, endTime, interval, includeMeta, perPage, offset}
	} else {
		// Get all points with pagination
		query = `
            SELECT id, vehicle_id, driver_id, latitude, longitude, 
                   altitude, speed, direction, accuracy, location_time,
                   CASE WHEN $4 = true THEN meta ELSE '{}'::jsonb END as meta
            FROM tbl_gps_location
            WHERE driver_id = $1 
            AND location_time BETWEEN $2 AND $3
            ORDER BY location_time
            LIMIT $5 OFFSET $6
        `
		args = []interface{}{driverID, startTime, endTime, includeMeta, perPage, offset}
	}

	// Execute query
	var locations []dto.GPSLocationDetails
	err = pgxscan.Select(context.Background(), db.DB, &locations, query, args...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error retrieving location history", err.Error()))
		return
	}

	// Calculate analytics
	var totalDistance float64
	var avgSpeed, maxSpeed float64

	// TODO:
	//if len(locations) > 0 {
	//	// Calculate distance and speeds only if we have points
	//	distanceQuery := `
	//        SELECT
	//            COALESCE(SUM(
	//                ST_DistanceSphere(
	//                    ST_MakePoint(longitude, latitude),
	//                    LAG(ST_MakePoint(longitude, latitude)) OVER (ORDER BY location_time)
	//                )
	//            ) / 1000, 0) as total_distance,
	//            COALESCE(AVG(speed), 0) as avg_speed,
	//            COALESCE(MAX(speed), 0) as max_speed
	//        FROM tbl_gps_location
	//        WHERE driver_id = $1
	//        AND location_time BETWEEN $2 AND $3
	//        AND speed IS NOT NULL
	//    `
	//	err = db.DB.QueryRow(
	//		context.Background(),
	//		distanceQuery,
	//		driverID, startTime, endTime,
	//	).Scan(&totalDistance, &avgSpeed, &maxSpeed)
	//
	//	if err != nil {
	//		// If error in analytics calculation, continue with zeros
	//		totalDistance, avgSpeed, maxSpeed = 0, 0, 0
	//	}
	//}

	// Create response
	response := gin.H{
		"driver_id":    driverID,
		"driver_name":  driverName,
		"total_points": totalCount,
		"points":       locations,
		"distance":     totalDistance,
		"avg_speed":    avgSpeed,
		"max_speed":    maxSpeed,
		"pagination": gin.H{
			"page":     page,
			"per_page": perPage,
			"total":    totalCount,
			"pages":    (totalCount + perPage - 1) / perPage,
		},
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Driver location history", response))
}

// GetTrip retrieves details about a specific trip
func GetTrip(ctx *gin.Context) {
	tripID := ctx.Param("trip_id")
	companyID := ctx.MustGet("companyID").(int)
	role := ctx.MustGet("role").(string)

	var trip dto.GPSTripDetails
	var query string

	if role == "admin" {
		query = `
            SELECT 
                t.*, 
                d.first_name || ' ' || d.last_name AS driver_name,
                v.numberplate AS vehicle_plate
            FROM tbl_gps_trip t
            JOIN tbl_driver d ON t.driver_id = d.id
            JOIN tbl_vehicle v ON t.vehicle_id = v.id
            WHERE t.uuid = $1
        `
	} else {
		query = `
            SELECT 
                t.*, 
                d.first_name || ' ' || d.last_name AS driver_name,
                v.numberplate AS vehicle_plate
            FROM tbl_gps_trip t
            JOIN tbl_driver d ON t.driver_id = d.id AND d.company_id = $2
            JOIN tbl_vehicle v ON t.vehicle_id = v.id AND v.company_id = $2
            WHERE t.uuid = $1
        `
	}

	err := pgxscan.Get(context.Background(), db.DB, &trip, query, tripID, companyID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Trip not found", err.Error()))
		return
	}

	//// TODO: Calculate trip duration if both start and end times exist
	//if trip.EndTime.Valid {
	//	duration := trip.EndTime.Time.Sub(trip.StartTime)
	//	trip.Duration = formatDuration(duration)
	//}

	// Fetch path points
	pointsQuery := `
        SELECT latitude, longitude
        FROM tbl_gps_location
        WHERE vehicle_id = $1
        AND location_time BETWEEN $2 AND $3
        ORDER BY location_time
    `

	var endTime time.Time
	if trip.EndTime.Valid {
		endTime = trip.EndTime.Time
	} else {
		endTime = time.Now()
	}

	type Point struct {
		Latitude  float64 `json:"latitude" db:"latitude"`
		Longitude float64 `json:"longitude" db:"longitude"`
	}

	var points []Point
	err = pgxscan.Select(context.Background(), db.DB, &points, pointsQuery, trip.VehicleID, trip.StartTime, endTime)
	if err != nil {
		// Continue even if path points retrieval fails
		points = []Point{}
	}

	// Convert to path format
	path := make([][]float64, len(points))
	for i, p := range points {
		path[i] = []float64{p.Longitude, p.Latitude}
	}

	response := trip
	additionalData := map[string]interface{}{
		"path": path,
	}

	// Merge the additional data with the response
	responseBytes, _ := json.Marshal(response)
	var responseMap map[string]interface{}
	json.Unmarshal(responseBytes, &responseMap)

	for k, v := range additionalData {
		responseMap[k] = v
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Trip details", responseMap))
}

// GetTripList retrieves a list of trips
func GetTripList(ctx *gin.Context) {
	vehicleIDStr := ctx.Query("vehicle_id")
	driverIDStr := ctx.Query("driver_id")
	startTimeStr := ctx.Query("start_time")
	endTimeStr := ctx.Query("end_time")
	statusStr := ctx.DefaultQuery("status", "")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
	offset := (page - 1) * perPage

	companyID := ctx.MustGet("companyID").(int)
	role := ctx.MustGet("role").(string)

	// Parse filter parameters
	var vehicleID, driverID int
	var startTime, endTime time.Time
	var err error

	if vehicleIDStr != "" {
		vehicleID, err = strconv.Atoi(vehicleIDStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid vehicle ID", err.Error()))
			return
		}
	}

	if driverIDStr != "" {
		driverID, err = strconv.Atoi(driverIDStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid driver ID", err.Error()))
			return
		}
	}

	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid start_time format", "Use ISO 8601 format"))
			return
		}
	} else {
		// Default to 30 days ago if not provided
		startTime = time.Now().AddDate(0, 0, -30)
	}

	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid end_time format", "Use ISO 8601 format"))
			return
		}
	} else {
		// Default to now if not provided
		endTime = time.Now()
	}

	// Build query
	var countQuery, dataQuery string
	var args []interface{}
	var whereClause []string

	// Always filter by company ID unless admin
	if role != "admin" {
		whereClause = append(whereClause, "v.company_id = $"+strconv.Itoa(len(args)+1))
		args = append(args, companyID)
	}

	// Add optional filters
	if vehicleID > 0 {
		whereClause = append(whereClause, "t.vehicle_id = $"+strconv.Itoa(len(args)+1))
		args = append(args, vehicleID)
	}

	if driverID > 0 {
		whereClause = append(whereClause, "t.driver_id = $"+strconv.Itoa(len(args)+1))
		args = append(args, driverID)
	}

	// Add time range filter
	whereClause = append(whereClause, "t.start_time >= $"+strconv.Itoa(len(args)+1))
	args = append(args, startTime)

	whereClause = append(whereClause, "t.start_time <= $"+strconv.Itoa(len(args)+1))
	args = append(args, endTime)

	// Add status filter if provided
	if statusStr != "" {
		whereClause = append(whereClause, "t.status = $"+strconv.Itoa(len(args)+1))
		args = append(args, statusStr)
	}

	// Build WHERE clause
	whereStatement := ""
	if len(whereClause) > 0 {
		whereStatement = "WHERE " + strings.Join(whereClause, " AND ")
	}

	// Build count query
	countQuery = `
        SELECT COUNT(*) 
        FROM tbl_gps_trip t
        JOIN tbl_vehicle v ON t.vehicle_id = v.id
        JOIN tbl_driver d ON t.driver_id = d.id
        ` + whereStatement

	// Build data query
	dataQuery = `
        SELECT 
            t.*,
            d.first_name || ' ' || d.last_name AS driver_name,
            v.numberplate AS vehicle_plate
        FROM tbl_gps_trip t
        JOIN tbl_vehicle v ON t.vehicle_id = v.id
        JOIN tbl_driver d ON t.driver_id = d.id
        ` + whereStatement + `
        ORDER BY t.start_time DESC
        LIMIT $` + strconv.Itoa(len(args)+1) + ` OFFSET $` + strconv.Itoa(len(args)+2)

	// Add pagination args
	args = append(args, perPage, offset)

	// Get count
	var totalCount int
	err = db.DB.QueryRow(context.Background(), countQuery, args[:len(args)-2]...).Scan(&totalCount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Database error", err.Error()))
		return
	}

	// Get data
	var trips []dto.GPSTripDetails
	err = pgxscan.Select(context.Background(), db.DB, &trips, dataQuery, args...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Database error", err.Error()))
		return
	}

	//// TODO: Calculate duration for each trip
	//for i := range trips {
	//	if trips[i].EndTime.Valid {
	//		duration := trips[i].EndTime.Time.Sub(trips[i].StartTime)
	//		trips[i].Duration = formatDuration(duration)
	//	}
	//}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Trip list", gin.H{
		"total":    totalCount,
		"page":     page,
		"per_page": perPage,
		"trips":    trips,
		"pages":    (totalCount + perPage - 1) / perPage,
	}))
}

// CreateTrip creates a new trip
func CreateTrip(ctx *gin.Context) {
	var trip dto.GPSTripCreate
	companyID := ctx.MustGet("companyID").(int)
	role := ctx.MustGet("role").(string)

	if err := ctx.ShouldBindJSON(&trip); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	// Verify vehicle and driver belong to the company
	if role != "admin" {
		var count int
		err := db.DB.QueryRow(
			context.Background(),
			`SELECT COUNT(*) FROM tbl_vehicle v
            JOIN tbl_driver d ON d.id = $1 AND d.company_id = $3
            WHERE v.id = $2 AND v.company_id = $3`,
			trip.DriverID, trip.VehicleID, companyID,
		).Scan(&count)
		if err != nil || count == 0 {
			ctx.JSON(http.StatusForbidden, utils.FormatErrorResponse("Permission denied", "Vehicle or driver does not belong to your company"))
			return
		}
	}

	// Set default values
	if trip.Status == "" {
		trip.Status = "active"
	}

	// Prepare optional fields
	var endTime, startLocation, endLocation, meta interface{} = nil, nil, nil, nil
	if trip.EndTime != nil {
		endTime = *trip.EndTime
	}
	if trip.StartLocation != nil {
		startLocation = trip.StartLocation
	}
	if trip.EndLocation != nil {
		endLocation = trip.EndLocation
	}
	if trip.Meta != nil {
		meta = trip.Meta
	} else {
		meta = []byte("{}")
	}

	// Insert trip
	query := `
        INSERT INTO tbl_gps_trip (
            vehicle_id, driver_id, start_time, end_time, 
            start_location, end_location, distance, avg_speed, 
            max_speed, status, meta
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
        ) RETURNING id, uuid
    `

	var tripID int
	var tripUUID string
	err := db.DB.QueryRow(
		context.Background(),
		query,
		trip.VehicleID, trip.DriverID, trip.StartTime, endTime,
		startLocation, endLocation, trip.Distance, trip.AvgSpeed,
		trip.MaxSpeed, trip.Status, meta,
	).Scan(&tripID, &tripUUID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Database error", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.FormatResponse("Trip created successfully", gin.H{
		"trip_id": tripID,
		"uuid":    tripUUID,
	}))
}

// UpdateTrip updates an existing trip
func UpdateTrip(ctx *gin.Context) {
	tripID := ctx.Param("trip_id")
	var tripUpdate dto.GPSTripUpdate
	companyID := ctx.MustGet("companyID").(int)
	role := ctx.MustGet("role").(string)

	if err := ctx.ShouldBindJSON(&tripUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	// Verify trip belongs to the company
	var vehicleID, driverID int
	var verifyQuery string
	if role == "admin" {
		verifyQuery = `
            SELECT vehicle_id, driver_id
            FROM tbl_gps_trip
            WHERE uuid = $1
        `
	} else {
		verifyQuery = `
            SELECT t.vehicle_id, t.driver_id
            FROM tbl_gps_trip t
            JOIN tbl_vehicle v ON t.vehicle_id = v.id AND v.company_id = $2
            JOIN tbl_driver d ON t.driver_id = d.id AND d.company_id = $2
            WHERE t.uuid = $1
        `
	}

	err := db.DB.QueryRow(
		context.Background(),
		verifyQuery,
		tripID, companyID,
	).Scan(&vehicleID, &driverID)

	if err != nil {
		ctx.JSON(http.StatusForbidden, utils.FormatErrorResponse("Permission denied", "Trip not found or does not belong to your company"))
		return
	}

	// Build update query
	var setClauses []string
	var args []interface{}
	args = append(args, tripID) // $1 is the trip UUID

	// Add fields to update
	if tripUpdate.EndTime != nil {
		setClauses = append(setClauses, "end_time = $"+strconv.Itoa(len(args)+1))
		args = append(args, *tripUpdate.EndTime)
	}

	if tripUpdate.Distance != nil {
		setClauses = append(setClauses, "distance = $"+strconv.Itoa(len(args)+1))
		args = append(args, *tripUpdate.Distance)
	}

	if tripUpdate.AvgSpeed != nil {
		setClauses = append(setClauses, "avg_speed = $"+strconv.Itoa(len(args)+1))
		args = append(args, *tripUpdate.AvgSpeed)
	}

	if tripUpdate.MaxSpeed != nil {
		setClauses = append(setClauses, "max_speed = $"+strconv.Itoa(len(args)+1))
		args = append(args, *tripUpdate.MaxSpeed)
	}

	if tripUpdate.Status != nil {
		setClauses = append(setClauses, "status = $"+strconv.Itoa(len(args)+1))
		args = append(args, *tripUpdate.Status)
	}

	if tripUpdate.EndLocation != nil {
		setClauses = append(setClauses, "end_location = $"+strconv.Itoa(len(args)+1))
		args = append(args, tripUpdate.EndLocation)
	}

	if tripUpdate.Meta != nil {
		setClauses = append(setClauses, "meta = $"+strconv.Itoa(len(args)+1))
		args = append(args, tripUpdate.Meta)
	}

	// Always update updated_at
	setClauses = append(setClauses, "updated_at = CURRENT_TIMESTAMP")

	// If no fields to update, return
	if len(setClauses) == 1 {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("No fields to update", "Request body is empty"))
		return
	}

	// Execute update
	query := fmt.Sprintf(`
        UPDATE tbl_gps_trip SET %s
        WHERE uuid = $1
        RETURNING id
    `, strings.Join(setClauses, ", "))

	var updatedID int
	err = db.DB.QueryRow(
		context.Background(),
		query,
		args...,
	).Scan(&updatedID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Database error", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Trip updated successfully", gin.H{
		"trip_id": updatedID,
		"uuid":    tripID,
	}))
}
