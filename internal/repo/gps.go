package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"texApi/pkg/utils"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	db "texApi/database"
	"texApi/internal/dto"
)

const (
	DefaultLimit    = 100
	DefaultOrderBy  = "id"
	DefaultOrderDir = "DESC"
)

func CreateTrip(input dto.StartTripInput) (int64, error) {
	ctx := context.Background()
	tx, err := db.DB.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var offerData struct {
		FromAddress *string `db:"from_address"`
		ToAddress   *string `db:"to_address"`
		FromCountry *string `db:"from_country"`
		ToCountry   *string `db:"to_country"`
		DriverID    int     `db:"driver_id"`
		VehicleID   int     `db:"vehicle_id"`
		DistanceKM  *int    `db:"distance"`
	}

	mainOfferID := input.Offers[0].OfferID
	err = pgxscan.Get(ctx, tx, &offerData,
		`SELECT from_address, to_address, from_country, to_country, 
		        driver_id, vehicle_id, distance
		 FROM tbl_offer WHERE id = $1 AND deleted = 0`,
		mainOfferID)
	if err != nil {
		return 0, fmt.Errorf("failed to get offer data: %w", err)
	}

	driverID := getIntValue(input.DriverID, offerData.DriverID)
	vehicleID := getIntValue(input.VehicleID, offerData.VehicleID)
	fromAddress := getStringValue(input.FromAddress, offerData.FromAddress)
	toAddress := getStringValue(input.ToAddress, offerData.ToAddress)
	fromCountry := getStringValue(input.FromCountry, offerData.FromCountry)
	toCountry := getStringValue(input.ToCountry, offerData.ToCountry)

	var distanceKM *float64
	if input.DistanceKM != nil {
		distanceKM = input.DistanceKM
	} else if offerData.DistanceKM != nil {
		distance := float64(*offerData.DistanceKM)
		distanceKM = &distance
	}

	startDate := input.StartDate
	if startDate == nil {
		now := time.Now()
		startDate = &now
	}

	var tripID int64
	err = pgxscan.Get(ctx, tx, &tripID,
		`INSERT INTO tbl_trip (
			driver_id, vehicle_id, from_address, to_address, 
			from_country, to_country, start_date, end_date,
			from_location, to_location, distance_km, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, 'active')
		RETURNING id`,
		driverID, vehicleID, fromAddress, toAddress,
		fromCountry, toCountry, startDate, input.EndDate,
		input.FromLocation, input.ToLocation, distanceKM)
	if err != nil {
		return 0, fmt.Errorf("failed to create trip: %w", err)
	}

	for _, offer := range input.Offers {
		_, err = tx.Exec(ctx,
			`INSERT INTO tbl_offer_trip (trip_id, offer_id, is_main, status)
			 VALUES ($1, $2, $3, 'active')`,
			tripID, offer.OfferID, offer.IsMain)
		if err != nil {
			return 0, fmt.Errorf("failed to link offer to trip: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return tripID, nil
}

func EndTrip(input dto.EndTripInput) error {
	ctx := context.Background()
	tx, err := db.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var logs []dto.GPSLog
	err = pgxscan.Select(ctx, tx, &logs,
		`SELECT * FROM tbl_gps_log 
		 WHERE trip_id = $1 
		 ORDER BY log_dt ASC`,
		input.ID)
	if err != nil {
		return fmt.Errorf("failed to get GPS logs: %w", err)
	}

	logsJSON, err := json.Marshal(logs)
	if err != nil {
		return fmt.Errorf("failed to marshal GPS logs: %w", err)
	}

	result, err := tx.Exec(ctx,
		`UPDATE tbl_trip 
		 SET status = 'completed', end_date = CURRENT_TIMESTAMP, 
		     gps_logs = $1, updated_at = CURRENT_TIMESTAMP
		 WHERE id = $2 AND driver_id = $3 AND deleted = 0`,
		string(logsJSON), input.ID, input.DriverID)
	if err != nil {
		return fmt.Errorf("failed to end trip: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("trip not found or access denied")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

type TripScan struct {
	ID              int64      `db:"id"`
	DriverID        int        `db:"driver_id"`
	VehicleID       int        `db:"vehicle_id"`
	FromAddress     *string    `db:"from_address"`
	ToAddress       *string    `db:"to_address"`
	FromCountry     *string    `db:"from_country"`
	ToCountry       *string    `db:"to_country"`
	StartDate       *time.Time `db:"start_date"`
	EndDate         *time.Time `db:"end_date"`
	FromLocationTxt *string    `db:"from_location_txt"` // ST_AsText result
	ToLocationTxt   *string    `db:"to_location_txt"`   // ST_AsText result
	DistanceKM      *float64   `db:"distance_km"`
	Status          string     `db:"status"`
	Meta            string     `db:"meta"`
	Meta2           string     `db:"meta2"`
	Meta3           string     `db:"meta3"`
	GPSLogs         string     `db:"gps_logs"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
	Deleted         int        `db:"deleted"`
}

func (ts *TripScan) ToTrip() dto.Trip {
	trip := dto.Trip{
		ID:          ts.ID,
		DriverID:    ts.DriverID,
		VehicleID:   ts.VehicleID,
		FromAddress: ts.FromAddress,
		ToAddress:   ts.ToAddress,
		FromCountry: ts.FromCountry,
		ToCountry:   ts.ToCountry,
		StartDate:   ts.StartDate,
		EndDate:     ts.EndDate,
		DistanceKM:  ts.DistanceKM,
		Status:      ts.Status,
		Meta:        ts.Meta,
		Meta2:       ts.Meta2,
		Meta3:       ts.Meta3,
		GPSLogs:     ts.GPSLogs,
		CreatedAt:   ts.CreatedAt,
		UpdatedAt:   ts.UpdatedAt,
		Deleted:     ts.Deleted,
	}

	if ts.FromLocationTxt != nil && *ts.FromLocationTxt != "" {
		var point dto.Point
		if err := point.Scan(*ts.FromLocationTxt); err == nil {
			trip.FromLocation = &point
		}
	}

	if ts.ToLocationTxt != nil && *ts.ToLocationTxt != "" {
		var point dto.Point
		if err := point.Scan(*ts.ToLocationTxt); err == nil {
			trip.ToLocation = &point
		}
	}

	return trip
}

type GPSLogScan struct {
	ID             int64     `db:"id"`
	CompanyID      *int      `db:"company_id"`
	VehicleID      int       `db:"vehicle_id"`
	DriverID       int       `db:"driver_id"`
	OfferID        *int      `db:"offer_id"`
	TripID         *int      `db:"trip_id"`
	BatteryLevel   *int      `db:"battery_level"`
	Speed          *float64  `db:"speed"`
	Heading        *float64  `db:"heading"`
	Accuracy       *float64  `db:"accuracy"`
	CoordinatesTxt string    `db:"coordinates_txt"` // ST_AsText result
	Status         string    `db:"status"`
	LogDt          time.Time `db:"log_dt"`
	CreatedAt      time.Time `db:"created_at"`
}

func (gs *GPSLogScan) ToGPSLog() dto.GPSLog {
	log := dto.GPSLog{
		ID:           gs.ID,
		CompanyID:    gs.CompanyID,
		VehicleID:    gs.VehicleID,
		DriverID:     gs.DriverID,
		OfferID:      gs.OfferID,
		TripID:       gs.TripID,
		BatteryLevel: gs.BatteryLevel,
		Speed:        gs.Speed,
		Heading:      gs.Heading,
		Accuracy:     gs.Accuracy,
		Status:       gs.Status,
		LogDt:        gs.LogDt,
		CreatedAt:    gs.CreatedAt,
	}

	if gs.CoordinatesTxt != "" {
		if err := log.Coordinates.Scan(gs.CoordinatesTxt); err != nil {
			// If parsing fails, set to zero coordinates
			log.Coordinates = dto.Point{Lat: 0, Lng: 0}
		}
	}

	return log
}

func GetTrips(query dto.TripQuery) ([]dto.Trip, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	conditions = append(conditions, "deleted = 0")

	offerFiltersQuery := dto.OfferFiltersQuery{
		OfferCompanyID:     query.OfferCompanyID,
		OfferExecCompanyID: query.OfferExecCompanyID,
		OfferDriverID:      query.OfferDriverID,
		OfferVehicleID:     query.OfferVehicleID,
		OfferFromCountryID: query.OfferFromCountryID,
		OfferToCountryID:   query.OfferToCountryID,
		OfferFromAddress:   query.OfferFromAddress,
		OfferToAddress:     query.OfferToAddress,
		OfferFromCountry:   query.OfferFromCountry,
		OfferToCountry:     query.OfferToCountry,
		OfferState:         query.OfferState,
		OfferRole:          query.OfferRole,
		OfferValidityStart: query.OfferValidityStart,
		OfferValidityEnd:   query.OfferValidityEnd,
		OfferDeliveryStart: query.OfferDeliveryStart,
		OfferDeliveryEnd:   query.OfferDeliveryEnd,
		OfferCostPerKmMin:  query.OfferCostPerKmMin,
		OfferCostPerKmMax:  query.OfferCostPerKmMax,
		OfferPriceMin:      query.OfferPriceMin,
		OfferPriceMax:      query.OfferPriceMax,
		OfferTotalPriceMin: query.OfferTotalPriceMin,
		OfferTotalPriceMax: query.OfferTotalPriceMax,
		Search:             query.Search,
	}

	if hasOfferFilters(offerFiltersQuery) {
		offerIDs, err := GetOfferIDsByParams(offerFiltersQuery)
		if err != nil {
			return nil, fmt.Errorf("failed to get offer IDs: %w", err)
		}

		if len(offerIDs) == 0 {
			return []dto.Trip{}, nil
		}

		placeholders := make([]string, len(offerIDs))
		for i, offerID := range offerIDs {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, offerID)
			argIndex++
		}

		conditions = append(conditions, fmt.Sprintf(`id IN (
            SELECT trip_id FROM tbl_offer_trip 
            WHERE offer_id IN (%s) AND deleted = 0
        )`, strings.Join(placeholders, ",")))
	}

	if query.DriverID != nil {
		conditions = append(conditions, fmt.Sprintf("driver_id = $%d", argIndex))
		args = append(args, *query.DriverID)
		argIndex++
	}

	if query.VehicleID != nil {
		conditions = append(conditions, fmt.Sprintf("vehicle_id = $%d", argIndex))
		args = append(args, *query.VehicleID)
		argIndex++
	}

	if query.FromAddress != nil {
		conditions = append(conditions, fmt.Sprintf("from_address ILIKE $%d", argIndex))
		args = append(args, "%"+*query.FromAddress+"%")
		argIndex++
	}

	if query.ToAddress != nil {
		conditions = append(conditions, fmt.Sprintf("to_address ILIKE $%d", argIndex))
		args = append(args, "%"+*query.ToAddress+"%")
		argIndex++
	}

	if query.FromCountry != nil {
		conditions = append(conditions, fmt.Sprintf("from_country ILIKE $%d", argIndex))
		args = append(args, "%"+*query.FromCountry+"%")
		argIndex++
	}

	if query.ToCountry != nil {
		conditions = append(conditions, fmt.Sprintf("to_country ILIKE $%d", argIndex))
		args = append(args, "%"+*query.ToCountry+"%")
		argIndex++
	}

	if query.StartDate != nil {
		conditions = append(conditions, fmt.Sprintf("start_date >= $%d", argIndex))
		args = append(args, *query.StartDate)
		argIndex++
	}

	if query.EndDate != nil {
		conditions = append(conditions, fmt.Sprintf("end_date <= $%d", argIndex))
		args = append(args, *query.EndDate)
		argIndex++
	}

	if query.DistanceKM != nil {
		conditions = append(conditions, fmt.Sprintf("distance_km >= $%d", argIndex))
		args = append(args, *query.DistanceKM)
		argIndex++
	}

	if query.TripOfferID != nil {
		conditions = append(conditions, fmt.Sprintf(`id IN (
            SELECT trip_id FROM tbl_offer_trip 
            WHERE offer_id = $%d AND deleted = 0
        )`, argIndex))
		args = append(args, *query.TripOfferID)
		argIndex++
	}

	if query.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *query.Status)
		argIndex++
	}

	if len(query.StatusIn) > 0 {
		placeholders := make([]string, len(query.StatusIn))
		for i, status := range query.StatusIn {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, status)
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("status IN (%s)", strings.Join(placeholders, ",")))
	}

	if query.StartDateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("start_date >= $%d", argIndex))
		args = append(args, *query.StartDateFrom)
		argIndex++
	}

	if query.StartDateTo != nil {
		conditions = append(conditions, fmt.Sprintf("start_date <= $%d", argIndex))
		args = append(args, *query.StartDateTo)
		argIndex++
	}

	if query.EndDateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("end_date >= $%d", argIndex))
		args = append(args, *query.EndDateFrom)
		argIndex++
	}

	if query.EndDateTo != nil {
		conditions = append(conditions, fmt.Sprintf("end_date <= $%d", argIndex))
		args = append(args, *query.EndDateTo)
		argIndex++
	}

	if query.CreatedAfter != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argIndex))
		args = append(args, *query.CreatedAfter)
		argIndex++
	}

	if query.CreatedBefore != nil {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argIndex))
		args = append(args, *query.CreatedBefore)
		argIndex++
	}

	if query.UpdatedAfter != nil {
		conditions = append(conditions, fmt.Sprintf("updated_at >= $%d", argIndex))
		args = append(args, *query.UpdatedAfter)
		argIndex++
	}

	if query.UpdatedBefore != nil {
		conditions = append(conditions, fmt.Sprintf("updated_at <= $%d", argIndex))
		args = append(args, *query.UpdatedBefore)
		argIndex++
	}

	if query.DistanceKMMin != nil {
		conditions = append(conditions, fmt.Sprintf("distance_km >= $%d", argIndex))
		args = append(args, *query.DistanceKMMin)
		argIndex++
	}

	if query.DistanceKMMax != nil {
		conditions = append(conditions, fmt.Sprintf("distance_km <= $%d", argIndex))
		args = append(args, *query.DistanceKMMax)
		argIndex++
	}

	if query.FromRegion != nil {
		conditions = append(conditions, fmt.Sprintf("from_address ILIKE $%d OR from_country ILIKE $%d", argIndex, argIndex))
		regionPattern := "%" + *query.FromRegion + "%"
		args = append(args, regionPattern)
		argIndex++
	}

	if query.ToRegion != nil {
		conditions = append(conditions, fmt.Sprintf("to_address ILIKE $%d OR to_country ILIKE $%d", argIndex, argIndex))
		regionPattern := "%" + *query.ToRegion + "%"
		args = append(args, regionPattern)
		argIndex++
	}

	if query.NearFromLat != nil && query.NearFromLng != nil && query.FromRadius != nil {
		conditions = append(conditions, fmt.Sprintf(
			"ST_DWithin(from_location, ST_SetSRID(ST_MakePoint($%d, $%d), 4326)::geography, $%d * 1000)",
			argIndex, argIndex+1, argIndex+2))
		args = append(args, *query.NearFromLng, *query.NearFromLat, *query.FromRadius)
		argIndex += 3
	}

	if query.NearToLat != nil && query.NearToLng != nil && query.ToRadius != nil {
		conditions = append(conditions, fmt.Sprintf(
			"ST_DWithin(to_location, ST_SetSRID(ST_MakePoint($%d, $%d), 4326)::geography, $%d * 1000)",
			argIndex, argIndex+1, argIndex+2))
		args = append(args, *query.NearToLng, *query.NearToLat, *query.ToRadius)
		argIndex += 3
	}

	if query.MetaContains != nil {
		conditions = append(conditions, fmt.Sprintf("meta ILIKE $%d", argIndex))
		args = append(args, "%"+*query.MetaContains+"%")
		argIndex++
	}

	if query.Meta2Contains != nil {
		conditions = append(conditions, fmt.Sprintf("meta2 ILIKE $%d", argIndex))
		args = append(args, "%"+*query.Meta2Contains+"%")
		argIndex++
	}

	if query.Meta3Contains != nil {
		conditions = append(conditions, fmt.Sprintf("meta3 ILIKE $%d", argIndex))
		args = append(args, "%"+*query.Meta3Contains+"%")
		argIndex++
	}

	if query.HasGPSLogs != nil {
		if *query.HasGPSLogs {
			conditions = append(conditions, "gps_logs != '{}'::jsonb AND gps_logs IS NOT NULL")
		} else {
			conditions = append(conditions, "(gps_logs = '{}'::jsonb OR gps_logs IS NULL)")
		}
	}

	if query.HasDriver != nil {
		if *query.HasDriver {
			conditions = append(conditions, "driver_id > 0")
		} else {
			conditions = append(conditions, "driver_id = 0")
		}
	}

	if query.HasVehicle != nil {
		if *query.HasVehicle {
			conditions = append(conditions, "vehicle_id > 0")
		} else {
			conditions = append(conditions, "vehicle_id = 0")
		}
	}

	if len(utils.SafeString(query.DriverIDs)) > 0 {
		placeholders := make([]string, len(utils.SafeString(query.DriverIDs)))
		for i, driverID := range utils.SafeString(query.DriverIDs) {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, driverID)
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("driver_id IN (%s)", strings.Join(placeholders, ",")))
	}

	if len(utils.SafeString(query.VehicleIDs)) > 0 {
		placeholders := make([]string, len(utils.SafeString(query.VehicleIDs)))
		for i, vehicleID := range utils.SafeString(query.VehicleIDs) {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, vehicleID)
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("vehicle_id IN (%s)", strings.Join(placeholders, ",")))
	}

	if len(utils.SafeString(query.TripIDs)) > 0 {
		placeholders := make([]string, len(utils.SafeString(query.TripIDs)))
		for i, tripID := range utils.SafeString(query.TripIDs) {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, tripID)
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("id IN (%s)", strings.Join(placeholders, ",")))
	}

	if query.ExcludeDriverID != nil {
		conditions = append(conditions, fmt.Sprintf("driver_id != $%d", argIndex))
		args = append(args, *query.ExcludeDriverID)
		argIndex++
	}

	if query.ExcludeVehicleID != nil {
		conditions = append(conditions, fmt.Sprintf("vehicle_id != $%d", argIndex))
		args = append(args, *query.ExcludeVehicleID)
		argIndex++
	}

	if len(utils.SafeString(query.ExcludeTripIDs)) > 0 {
		placeholders := make([]string, len(utils.SafeString(query.ExcludeTripIDs)))
		for i, tripID := range utils.SafeString(query.ExcludeTripIDs) {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, tripID)
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("id NOT IN (%s)", strings.Join(placeholders, ",")))
	}

	if query.Search != nil {
		searchPattern := "%" + *query.Search + "%"
		conditions = append(conditions, fmt.Sprintf(`(
            from_address ILIKE $%d OR 
            to_address ILIKE $%d OR 
            from_country ILIKE $%d OR 
            to_country ILIKE $%d OR
            meta ILIKE $%d OR
            meta2 ILIKE $%d OR
            meta3 ILIKE $%d
        )`, argIndex, argIndex, argIndex, argIndex, argIndex, argIndex, argIndex))
		args = append(args, searchPattern)
		argIndex++
	}

	limit := query.Limit
	if limit == 0 {
		limit = DefaultLimit
	}

	orderBy := utils.SafeString(query.OrderBy)
	if orderBy == "" {
		orderBy = DefaultOrderBy
	}

	orderDir := utils.SafeString(query.OrderDir)
	if orderDir == "" {
		orderDir = DefaultOrderDir
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	queryStr := fmt.Sprintf(`
        SELECT id, driver_id, vehicle_id, from_address, to_address, 
               from_country, to_country, start_date, end_date,
               ST_AsText(from_location) as from_location_txt,
               ST_AsText(to_location) as to_location_txt,
               distance_km, status, meta, meta2, meta3, gps_logs,
               created_at, updated_at, deleted
        FROM tbl_trip 
        %s 
        ORDER BY %s %s 
        LIMIT $%d OFFSET $%d`,
		whereClause, orderBy, orderDir, argIndex, argIndex+1)

	args = append(args, limit, query.Offset)

	var tripScans []TripScan
	err := pgxscan.Select(context.Background(), db.DB, &tripScans, queryStr, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get trips: %w", err)
	}

	trips := make([]dto.Trip, len(tripScans))
	for i, scan := range tripScans {
		trips[i] = scan.ToTrip()
	}

	return trips, nil
}

func GetGPSLogs(query dto.GPSLogQuery) ([]dto.GPSLog, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if query.TripID != nil {
		conditions = append(conditions, fmt.Sprintf("trip_id = $%d", argIndex))
		args = append(args, *query.TripID)
		argIndex++
	}

	if query.CompanyID != nil {
		conditions = append(conditions, fmt.Sprintf("company_id = $%d", argIndex))
		args = append(args, *query.CompanyID)
		argIndex++
	}

	if query.OfferID != nil {
		conditions = append(conditions, fmt.Sprintf("offer_id = $%d", argIndex))
		args = append(args, *query.OfferID)
		argIndex++
	}

	if query.DriverID != nil {
		conditions = append(conditions, fmt.Sprintf("driver_id = $%d", argIndex))
		args = append(args, *query.DriverID)
		argIndex++
	}

	if query.VehicleID != nil {
		conditions = append(conditions, fmt.Sprintf("vehicle_id = $%d", argIndex))
		args = append(args, *query.VehicleID)
		argIndex++
	}

	if query.From != nil {
		conditions = append(conditions, fmt.Sprintf("log_dt >= $%d", argIndex))
		args = append(args, *query.From)
		argIndex++
	}

	if query.To != nil {
		conditions = append(conditions, fmt.Sprintf("log_dt <= $%d", argIndex))
		args = append(args, *query.To)
		argIndex++
	}

	if query.TripOfferID != nil {
		conditions = append(conditions, fmt.Sprintf(`trip_id IN (
			SELECT trip_id FROM tbl_offer_trip 
			WHERE offer_id = $%d AND deleted = 0
		)`, argIndex))
		args = append(args, *query.TripOfferID)
		argIndex++
	}

	limit := query.Limit
	if limit == 0 {
		limit = DefaultLimit
	}

	orderBy := utils.SafeString(query.OrderBy)
	if orderBy == "" {
		orderBy = DefaultOrderBy
	}

	orderDir := utils.SafeString(query.OrderDir)
	if orderDir == "" {
		orderDir = DefaultOrderDir
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	queryStr := fmt.Sprintf(`
		SELECT id, company_id, vehicle_id, driver_id, offer_id, trip_id,
		       battery_level, speed, heading, accuracy,
		       ST_AsText(coordinates) as coordinates_txt,
		       status, log_dt, created_at
		FROM tbl_gps_log 
		%s 
		ORDER BY %s %s 
		LIMIT $%d OFFSET $%d`,
		whereClause, orderBy, orderDir, argIndex, argIndex+1)

	args = append(args, limit, query.Offset)

	var logScans []GPSLogScan
	err := pgxscan.Select(context.Background(), db.DB, &logScans, queryStr, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get GPS logs: %w", err)
	}

	logs := make([]dto.GPSLog, len(logScans))
	for i, scan := range logScans {
		logs[i] = scan.ToGPSLog()
	}

	return logs, nil
}

func GetLastPositions(query dto.PositionQuery) ([]dto.GPSLog, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	offerFiltersQuery := dto.OfferFiltersQuery{
		OfferCompanyID:     query.OfferCompanyID,
		OfferExecCompanyID: query.OfferExecCompanyID,
		OfferDriverID:      query.OfferDriverID,
		OfferVehicleID:     query.OfferVehicleID,
		OfferFromCountryID: query.OfferFromCountryID,
		OfferToCountryID:   query.OfferToCountryID,
		OfferFromAddress:   query.OfferFromAddress,
		OfferToAddress:     query.OfferToAddress,
		OfferFromCountry:   query.OfferFromCountry,
		OfferToCountry:     query.OfferToCountry,
		OfferState:         query.OfferState,
		OfferRole:          query.OfferRole,
		OfferValidityStart: query.OfferValidityStart,
		OfferValidityEnd:   query.OfferValidityEnd,
		OfferDeliveryStart: query.OfferDeliveryStart,
		OfferDeliveryEnd:   query.OfferDeliveryEnd,
		OfferCostPerKmMin:  query.OfferCostPerKmMin,
		OfferCostPerKmMax:  query.OfferCostPerKmMax,
		OfferPriceMin:      query.OfferPriceMin,
		OfferPriceMax:      query.OfferPriceMax,
		OfferTotalPriceMin: query.OfferTotalPriceMin,
		OfferTotalPriceMax: query.OfferTotalPriceMax,
		Search:             query.Search,
	}

	if hasOfferFilters(offerFiltersQuery) {
		offerIDs, err := GetOfferIDsByParams(offerFiltersQuery)
		if err != nil {
			return nil, fmt.Errorf("failed to get offer IDs: %w", err)
		}

		if len(offerIDs) == 0 {
			return []dto.GPSLog{}, nil
		}

		placeholders := make([]string, len(offerIDs))
		for i, offerID := range offerIDs {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, offerID)
			argIndex++
		}

		conditions = append(conditions, fmt.Sprintf(`id IN (
            SELECT trip_id FROM tbl_offer_trip 
            WHERE offer_id IN (%s) AND deleted = 0
        )`, strings.Join(placeholders, ",")))
	}

	if len(query.TripIDs) > 0 {
		placeholders := make([]string, len(query.TripIDs))
		for i, tripID := range query.TripIDs {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, tripID)
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("trip_id IN (%s)", strings.Join(placeholders, ",")))
	}

	if len(query.CompanyIDs) > 0 {
		placeholders := make([]string, len(query.CompanyIDs))
		for i, companyID := range query.CompanyIDs {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, companyID)
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("company_id IN (%s)", strings.Join(placeholders, ",")))
	}

	if len(query.OfferIDs) > 0 {
		placeholders := make([]string, len(query.OfferIDs))
		for i, offerID := range query.OfferIDs {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, offerID)
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("offer_id IN (%s)", strings.Join(placeholders, ",")))
	}

	if len(query.DriverIDs) > 0 {
		placeholders := make([]string, len(query.DriverIDs))
		for i, driverID := range query.DriverIDs {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, driverID)
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("driver_id IN (%s)", strings.Join(placeholders, ",")))
	}

	if len(query.VehicleIDs) > 0 {
		placeholders := make([]string, len(query.VehicleIDs))
		for i, vehicleID := range query.VehicleIDs {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, vehicleID)
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("vehicle_id IN (%s)", strings.Join(placeholders, ",")))
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	queryStr := fmt.Sprintf(`
        SELECT DISTINCT ON (COALESCE(trip_id, 0), COALESCE(driver_id, 0), COALESCE(vehicle_id, 0)) 
               id, company_id, vehicle_id, driver_id, offer_id, trip_id,
               battery_level, speed, heading, accuracy,
               ST_AsText(coordinates) as coordinates_txt,
               status, log_dt, created_at
        FROM tbl_gps_log 
        %s 
        ORDER BY COALESCE(trip_id, 0), COALESCE(driver_id, 0), COALESCE(vehicle_id, 0), log_dt DESC`,
		whereClause)

	var logScans []GPSLogScan
	err := pgxscan.Select(context.Background(), db.DB, &logScans, queryStr, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get last positions: %w", err)
	}

	// Convert scanned results to dto.GPSLog
	logs := make([]dto.GPSLog, len(logScans))
	for i, scan := range logScans {
		logs[i] = scan.ToGPSLog()
	}

	return logs, nil
}

func CreateGPSLogs(logs []dto.GPSLogInput) error {
	if len(logs) == 0 {
		return nil
	}

	ctx := context.Background()
	tx, err := db.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	for _, log := range logs {
		_, err = tx.Exec(ctx,
			`INSERT INTO tbl_gps_log (
				company_id, vehicle_id, driver_id, offer_id, trip_id,
				battery_level, speed, heading, accuracy, coordinates,
				status, log_dt, created_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, 'active', $11, CURRENT_TIMESTAMP)`,
			log.CompanyID, log.VehicleID, log.DriverID, log.OfferID, log.TripID,
			log.BatteryLevel, log.Speed, log.Heading, log.Accuracy, log.Coordinates,
			log.LogDt)
		if err != nil {
			return fmt.Errorf("failed to create GPS log: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func getIntValue(input *int, fallback int) int {
	if input != nil {
		return *input
	}
	return fallback
}

func getStringValue(input *string, fallback *string) *string {
	if input != nil {
		return input
	}
	return fallback
}

type TripDetailedScan struct {
	ID              int64            `db:"id"`
	DriverID        int              `db:"driver_id"`
	VehicleID       int              `db:"vehicle_id"`
	FromAddress     *string          `db:"from_address"`
	ToAddress       *string          `db:"to_address"`
	FromCountry     *string          `db:"from_country"`
	ToCountry       *string          `db:"to_country"`
	StartDate       *time.Time       `db:"start_date"`
	EndDate         *time.Time       `db:"end_date"`
	FromLocationTxt *string          `db:"from_location_txt"`
	ToLocationTxt   *string          `db:"to_location_txt"`
	DistanceKM      *float64         `db:"distance_km"`
	Status          string           `db:"status"`
	Meta            string           `db:"meta"`
	Meta2           string           `db:"meta2"`
	Meta3           string           `db:"meta3"`
	GPSLogs         string           `db:"gps_logs"`
	CreatedAt       time.Time        `db:"created_at"`
	UpdatedAt       time.Time        `db:"updated_at"`
	Deleted         int              `db:"deleted"`
	TotalCount      int              `db:"total_count"`
	Driver          *json.RawMessage `db:"driver"`
	Vehicle         *json.RawMessage `db:"vehicle"`
	Offers          *json.RawMessage `db:"offers"`
}

func (ts *TripDetailedScan) ToTripDetailed() dto.TripDetailed {
	trip := dto.TripDetailed{
		ID:          ts.ID,
		DriverID:    ts.DriverID,
		VehicleID:   ts.VehicleID,
		FromAddress: ts.FromAddress,
		ToAddress:   ts.ToAddress,
		FromCountry: ts.FromCountry,
		ToCountry:   ts.ToCountry,
		StartDate:   ts.StartDate,
		EndDate:     ts.EndDate,
		DistanceKM:  ts.DistanceKM,
		Status:      ts.Status,
		Meta:        ts.Meta,
		Meta2:       ts.Meta2,
		Meta3:       ts.Meta3,
		GPSLogs:     ts.GPSLogs,
		CreatedAt:   ts.CreatedAt,
		UpdatedAt:   ts.UpdatedAt,
		Deleted:     ts.Deleted,
		TotalCount:  ts.TotalCount,
		Driver:      ts.Driver,
		Vehicle:     ts.Vehicle,
		Offers:      ts.Offers,
	}

	if ts.FromLocationTxt != nil && *ts.FromLocationTxt != "" {
		var point dto.Point
		if err := point.Scan(*ts.FromLocationTxt); err == nil {
			trip.FromLocation = &point
		}
	}

	if ts.ToLocationTxt != nil && *ts.ToLocationTxt != "" {
		var point dto.Point
		if err := point.Scan(*ts.ToLocationTxt); err == nil {
			trip.ToLocation = &point
		}
	}

	return trip
}

func GetTripsDetailed(query dto.TripQuery) ([]dto.TripDetailed, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	conditions = append(conditions, "t.deleted = 0")

	offerFiltersQuery := dto.OfferFiltersQuery{
		OfferCompanyID:     query.OfferCompanyID,
		OfferExecCompanyID: query.OfferExecCompanyID,
		OfferDriverID:      query.OfferDriverID,
		OfferVehicleID:     query.OfferVehicleID,
		OfferFromCountryID: query.OfferFromCountryID,
		OfferToCountryID:   query.OfferToCountryID,
		OfferFromAddress:   query.OfferFromAddress,
		OfferToAddress:     query.OfferToAddress,
		OfferFromCountry:   query.OfferFromCountry,
		OfferToCountry:     query.OfferToCountry,
		OfferState:         query.OfferState,
		OfferRole:          query.OfferRole,
		OfferValidityStart: query.OfferValidityStart,
		OfferValidityEnd:   query.OfferValidityEnd,
		OfferDeliveryStart: query.OfferDeliveryStart,
		OfferDeliveryEnd:   query.OfferDeliveryEnd,
		OfferCostPerKmMin:  query.OfferCostPerKmMin,
		OfferCostPerKmMax:  query.OfferCostPerKmMax,
		OfferPriceMin:      query.OfferPriceMin,
		OfferPriceMax:      query.OfferPriceMax,
		OfferTotalPriceMin: query.OfferTotalPriceMin,
		OfferTotalPriceMax: query.OfferTotalPriceMax,
		Search:             query.Search,
	}

	if hasOfferFilters(offerFiltersQuery) {
		offerIDs, err := GetOfferIDsByParams(offerFiltersQuery)
		if err != nil {
			return nil, fmt.Errorf("failed to get offer IDs: %w", err)
		}

		if len(offerIDs) == 0 {
			return []dto.TripDetailed{}, nil
		}

		placeholders := make([]string, len(offerIDs))
		for i, offerID := range offerIDs {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, offerID)
			argIndex++
		}

		conditions = append(conditions, fmt.Sprintf(`t.id IN (
			SELECT trip_id FROM tbl_offer_trip 
			WHERE offer_id IN (%s) AND deleted = 0
		)`, strings.Join(placeholders, ",")))
	}

	if query.DriverID != nil {
		conditions = append(conditions, fmt.Sprintf("t.driver_id = $%d", argIndex))
		args = append(args, *query.DriverID)
		argIndex++
	}

	if query.VehicleID != nil {
		conditions = append(conditions, fmt.Sprintf("t.vehicle_id = $%d", argIndex))
		args = append(args, *query.VehicleID)
		argIndex++
	}

	if query.FromAddress != nil {
		conditions = append(conditions, fmt.Sprintf("t.from_address ILIKE $%d", argIndex))
		args = append(args, "%"+*query.FromAddress+"%")
		argIndex++
	}

	if query.ToAddress != nil {
		conditions = append(conditions, fmt.Sprintf("t.to_address ILIKE $%d", argIndex))
		args = append(args, "%"+*query.ToAddress+"%")
		argIndex++
	}

	if query.FromCountry != nil {
		conditions = append(conditions, fmt.Sprintf("t.from_country ILIKE $%d", argIndex))
		args = append(args, "%"+*query.FromCountry+"%")
		argIndex++
	}

	if query.ToCountry != nil {
		conditions = append(conditions, fmt.Sprintf("t.to_country ILIKE $%d", argIndex))
		args = append(args, "%"+*query.ToCountry+"%")
		argIndex++
	}

	if query.StartDate != nil {
		conditions = append(conditions, fmt.Sprintf("t.start_date >= $%d", argIndex))
		args = append(args, *query.StartDate)
		argIndex++
	}

	if query.EndDate != nil {
		conditions = append(conditions, fmt.Sprintf("t.end_date <= $%d", argIndex))
		args = append(args, *query.EndDate)
		argIndex++
	}

	if query.DistanceKM != nil {
		conditions = append(conditions, fmt.Sprintf("t.distance_km >= $%d", argIndex))
		args = append(args, *query.DistanceKM)
		argIndex++
	}

	if query.TripOfferID != nil {
		conditions = append(conditions, fmt.Sprintf(`t.id IN (
			SELECT trip_id FROM tbl_offer_trip 
			WHERE offer_id = $%d AND deleted = 0
		)`, argIndex))
		args = append(args, *query.TripOfferID)
		argIndex++
	}

	if query.Status != nil {
		conditions = append(conditions, fmt.Sprintf("t.status = $%d", argIndex))
		args = append(args, *query.Status)
		argIndex++
	}

	if len(query.StatusIn) > 0 {
		placeholders := make([]string, len(query.StatusIn))
		for i, status := range query.StatusIn {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, status)
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("t.status IN (%s)", strings.Join(placeholders, ",")))
	}

	if query.StartDateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("t.start_date >= $%d", argIndex))
		args = append(args, *query.StartDateFrom)
		argIndex++
	}

	if query.StartDateTo != nil {
		conditions = append(conditions, fmt.Sprintf("t.start_date <= $%d", argIndex))
		args = append(args, *query.StartDateTo)
		argIndex++
	}

	if query.EndDateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("t.end_date >= $%d", argIndex))
		args = append(args, *query.EndDateFrom)
		argIndex++
	}

	if query.EndDateTo != nil {
		conditions = append(conditions, fmt.Sprintf("t.end_date <= $%d", argIndex))
		args = append(args, *query.EndDateTo)
		argIndex++
	}

	if query.CreatedAfter != nil {
		conditions = append(conditions, fmt.Sprintf("t.created_at >= $%d", argIndex))
		args = append(args, *query.CreatedAfter)
		argIndex++
	}

	if query.CreatedBefore != nil {
		conditions = append(conditions, fmt.Sprintf("t.created_at <= $%d", argIndex))
		args = append(args, *query.CreatedBefore)
		argIndex++
	}

	if query.UpdatedAfter != nil {
		conditions = append(conditions, fmt.Sprintf("t.updated_at >= $%d", argIndex))
		args = append(args, *query.UpdatedAfter)
		argIndex++
	}

	if query.UpdatedBefore != nil {
		conditions = append(conditions, fmt.Sprintf("t.updated_at <= $%d", argIndex))
		args = append(args, *query.UpdatedBefore)
		argIndex++
	}

	if query.DistanceKMMin != nil {
		conditions = append(conditions, fmt.Sprintf("t.distance_km >= $%d", argIndex))
		args = append(args, *query.DistanceKMMin)
		argIndex++
	}

	if query.DistanceKMMax != nil {
		conditions = append(conditions, fmt.Sprintf("t.distance_km <= $%d", argIndex))
		args = append(args, *query.DistanceKMMax)
		argIndex++
	}

	if query.FromRegion != nil {
		conditions = append(conditions, fmt.Sprintf("t.from_address ILIKE $%d OR t.from_country ILIKE $%d", argIndex, argIndex))
		regionPattern := "%" + *query.FromRegion + "%"
		args = append(args, regionPattern)
		argIndex++
	}

	if query.ToRegion != nil {
		conditions = append(conditions, fmt.Sprintf("t.to_address ILIKE $%d OR t.to_country ILIKE $%d", argIndex, argIndex))
		regionPattern := "%" + *query.ToRegion + "%"
		args = append(args, regionPattern)
		argIndex++
	}

	if query.NearFromLat != nil && query.NearFromLng != nil && query.FromRadius != nil {
		conditions = append(conditions, fmt.Sprintf(
			"ST_DWithin(t.from_location, ST_SetSRID(ST_MakePoint($%d, $%d), 4326)::geography, $%d * 1000)",
			argIndex, argIndex+1, argIndex+2))
		args = append(args, *query.NearFromLng, *query.NearFromLat, *query.FromRadius)
		argIndex += 3
	}

	if query.NearToLat != nil && query.NearToLng != nil && query.ToRadius != nil {
		conditions = append(conditions, fmt.Sprintf(
			"ST_DWithin(t.to_location, ST_SetSRID(ST_MakePoint($%d, $%d), 4326)::geography, $%d * 1000)",
			argIndex, argIndex+1, argIndex+2))
		args = append(args, *query.NearToLng, *query.NearToLat, *query.ToRadius)
		argIndex += 3
	}

	if query.MetaContains != nil {
		conditions = append(conditions, fmt.Sprintf("t.meta ILIKE $%d", argIndex))
		args = append(args, "%"+*query.MetaContains+"%")
		argIndex++
	}

	if query.Meta2Contains != nil {
		conditions = append(conditions, fmt.Sprintf("t.meta2 ILIKE $%d", argIndex))
		args = append(args, "%"+*query.Meta2Contains+"%")
		argIndex++
	}

	if query.Meta3Contains != nil {
		conditions = append(conditions, fmt.Sprintf("t.meta3 ILIKE $%d", argIndex))
		args = append(args, "%"+*query.Meta3Contains+"%")
		argIndex++
	}

	if query.HasGPSLogs != nil {
		if *query.HasGPSLogs {
			conditions = append(conditions, "t.gps_logs != '{}'::jsonb AND t.gps_logs IS NOT NULL")
		} else {
			conditions = append(conditions, "(t.gps_logs = '{}'::jsonb OR t.gps_logs IS NULL)")
		}
	}

	if query.HasDriver != nil {
		if *query.HasDriver {
			conditions = append(conditions, "t.driver_id > 0")
		} else {
			conditions = append(conditions, "t.driver_id = 0")
		}
	}

	if query.HasVehicle != nil {
		if *query.HasVehicle {
			conditions = append(conditions, "t.vehicle_id > 0")
		} else {
			conditions = append(conditions, "t.vehicle_id = 0")
		}
	}

	if len(utils.SafeString(query.DriverIDs)) > 0 {
		placeholders := make([]string, len(utils.SafeString(query.DriverIDs)))
		for i, driverID := range utils.SafeString(query.DriverIDs) {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, driverID)
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("t.driver_id IN (%s)", strings.Join(placeholders, ",")))
	}

	if len(utils.SafeString(query.VehicleIDs)) > 0 {
		placeholders := make([]string, len(utils.SafeString(query.VehicleIDs)))
		for i, vehicleID := range utils.SafeString(query.VehicleIDs) {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, vehicleID)
			argIndex++
		}
		fmt.Println(placeholders)
		conditions = append(conditions, fmt.Sprintf("t.vehicle_id IN (%s)", strings.Join(placeholders, ",")))
		fmt.Println(conditions)

	}

	if len(utils.SafeString(query.TripIDs)) > 0 {
		placeholders := make([]string, len(utils.SafeString(query.TripIDs)))
		for i, tripID := range utils.SafeString(query.TripIDs) {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, tripID)
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("t.id IN (%s)", strings.Join(placeholders, ",")))
	}

	if query.ExcludeDriverID != nil {
		conditions = append(conditions, fmt.Sprintf("t.driver_id != $%d", argIndex))
		args = append(args, *query.ExcludeDriverID)
		argIndex++
	}

	if query.ExcludeVehicleID != nil {
		conditions = append(conditions, fmt.Sprintf("t.vehicle_id != $%d", argIndex))
		args = append(args, *query.ExcludeVehicleID)
		argIndex++
	}

	if len(utils.SafeString(query.ExcludeTripIDs)) > 0 {
		placeholders := make([]string, len(utils.SafeString(query.ExcludeTripIDs)))
		for i, tripID := range utils.SafeString(query.ExcludeTripIDs) {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, tripID)
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("t.id NOT IN (%s)", strings.Join(placeholders, ",")))
	}

	if query.Search != nil {
		searchPattern := "%" + *query.Search + "%"
		conditions = append(conditions, fmt.Sprintf(`(
			t.from_address ILIKE $%d OR 
			t.to_address ILIKE $%d OR 
			t.from_country ILIKE $%d OR 
			t.to_country ILIKE $%d OR
			t.meta ILIKE $%d OR
			t.meta2 ILIKE $%d OR
			t.meta3 ILIKE $%d
		)`, argIndex, argIndex, argIndex, argIndex, argIndex, argIndex, argIndex))
		args = append(args, searchPattern)
		argIndex++
	}

	limit := query.Limit
	if limit == 0 {
		limit = DefaultLimit
	}

	orderBy := utils.SafeString(query.OrderBy)
	if orderBy == "" {
		orderBy = DefaultOrderBy
	}

	orderDir := utils.SafeString(query.OrderDir)
	if orderDir == "" {
		orderDir = DefaultOrderDir
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	queryStr := fmt.Sprintf(`
        SELECT 
            t.id, t.driver_id, t.vehicle_id, t.from_address, t.to_address, 
            t.from_country, t.to_country, t.start_date, t.end_date,
            ST_AsText(t.from_location) as from_location_txt,
            ST_AsText(t.to_location) as to_location_txt,
            t.distance_km, t.status, t.meta, t.meta2, t.meta3, t.gps_logs,
            t.created_at, t.updated_at, t.deleted,
            COUNT(*) OVER() as total_count,
            CASE 
                WHEN t.driver_id > 0 THEN 
                    json_build_object(
                        'id', d.id,
                        'uuid', d.uuid,
                        'company_id', d.company_id,
                        'first_name', d.first_name,
                        'last_name', d.last_name,
                        'patronymic_name', d.patronymic_name,
                        'phone', d.phone,
                        'email', d.email,
                        'featured', d.featured,
                        'rating', d.rating,
                        'partner', d.partner,
                        'successful_ops', d.successful_ops,
                        'image_url', d.image_url,
                        'view_count', d.view_count,
                        'meta', d.meta,
                        'meta2', d.meta2,
                        'meta3', d.meta3,
                        'available', d.available,
                        'block_reason', d.block_reason
                    )
                ELSE NULL
            END as driver,
            CASE 
                WHEN t.vehicle_id > 0 THEN 
                    json_build_object(
                        'id', v.id,
                        'uuid', v.uuid,
                        'company_id', v.company_id,
                        'vehicle_type_id', v.vehicle_type_id,
                        'vehicle_brand_id', v.vehicle_brand_id,
                        'vehicle_model_id', v.vehicle_model_id,
                        'year_of_issue', v.year_of_issue,
                        'mileage', v.mileage,
                        'numberplate', v.numberplate,
                        'trailer_numberplate', v.trailer_numberplate,
                        'gps', v.gps,
                        'photo1_url', v.photo1_url,
                        'photo2_url', v.photo2_url,
                        'photo3_url', v.photo3_url,
                        'docs1_url', v.docs1_url,
                        'docs2_url', v.docs2_url,
                        'docs3_url', v.docs3_url,
                        'view_count', v.view_count,
                        'meta', v.meta,
                        'meta2', v.meta2,
                        'meta3', v.meta3,
                        'available', v.available
                    )
                ELSE NULL
            END as vehicle,
            COALESCE((
                SELECT json_agg(
                    json_build_object(
                        'id', o.id,
                        'uuid', o.uuid,
                        'user_id', o.user_id,
                        'company_id', o.company_id,
                        'exec_company_id', o.exec_company_id,
                        'driver_id', o.driver_id,
                        'vehicle_id', o.vehicle_id,
                        'trailer_id', o.trailer_id,
                        'vehicle_type_id', o.vehicle_type_id,
                        'cargo_id', o.cargo_id,
                        'packaging_type_id', o.packaging_type_id,
                        'offer_state', o.offer_state,
                        'offer_role', o.offer_role,
                        'cost_per_km', o.cost_per_km,
                        'currency', o.currency,
                        'from_country_id', o.from_country_id,
                        'from_city_id', o.from_city_id,
                        'to_country_id', o.to_country_id,
                        'to_city_id', o.to_city_id,
                        'distance', o.distance,
                        'from_country', o.from_country,
                        'from_region', o.from_region,
                        'to_country', o.to_country,
                        'to_region', o.to_region,
                        'from_address', o.from_address,
                        'to_address', o.to_address,
                        'map_url', o.map_url,
                        'sender_contact', o.sender_contact,
                        'recipient_contact', o.recipient_contact,
                        'deliver_contact', o.deliver_contact,
                        'view_count', o.view_count,
                        'validity_start', o.validity_start,
                        'validity_end', o.validity_end,
                        'delivery_start', o.delivery_start,
                        'delivery_end', o.delivery_end,
                        'note', o.note,
                        'tax', o.tax,
                        'tax_price', o.tax_price,
                        'trade', o.trade,
                        'discount', o.discount,
                        'payment_method', o.payment_method,
                        'payment_term', o.payment_term,
                        'meta', o.meta,
                        'meta2', o.meta2,
                        'meta3', o.meta3,
                        'featured', o.featured,
                        'partner', o.partner,
                        'is_main', ot.is_main
                    )
                )
                FROM tbl_offer_trip ot
                JOIN tbl_offer o ON ot.offer_id = o.id
                WHERE ot.trip_id = t.id AND ot.deleted = 0 AND o.deleted = 0
            ), '[]') as offers
        FROM tbl_trip t
        LEFT JOIN tbl_driver d ON t.driver_id = d.id AND d.deleted = 0
        LEFT JOIN tbl_vehicle v ON t.vehicle_id = v.id AND v.deleted = 0
        %s 
        ORDER BY t.%s %s 
        LIMIT $%d OFFSET $%d`,
		whereClause, orderBy, orderDir, argIndex, argIndex+1)

	args = append(args, limit, query.Offset)

	var tripScans []TripDetailedScan
	err := pgxscan.Select(context.Background(), db.DB, &tripScans, queryStr, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get detailed trips: %w", err)
	}

	trips := make([]dto.TripDetailed, len(tripScans))
	for i, scan := range tripScans {
		trips[i] = scan.ToTripDetailed()
	}

	return trips, nil
}

func GetOfferIDsByParams(query dto.OfferFiltersQuery) ([]int, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	conditions = append(conditions, "o.deleted = 0")

	if query.Search != nil && strings.TrimSpace(*query.Search) != "" {
		searchTerm := "%" + strings.TrimSpace(*query.Search) + "%"

		searchConditions := []string{
			fmt.Sprintf("o.note ILIKE $%d", argIndex),
			fmt.Sprintf("o.meta ILIKE $%d", argIndex),
			fmt.Sprintf("o.from_address ILIKE $%d", argIndex),
			fmt.Sprintf("o.to_address ILIKE $%d", argIndex),
			fmt.Sprintf("o.from_country ILIKE $%d", argIndex),
			fmt.Sprintf("o.to_country ILIKE $%d", argIndex),
			fmt.Sprintf("o.from_region ILIKE $%d", argIndex),
			fmt.Sprintf("o.to_region ILIKE $%d", argIndex),
			fmt.Sprintf("o.sender_contact ILIKE $%d", argIndex),
			fmt.Sprintf("o.recipient_contact ILIKE $%d", argIndex),
			fmt.Sprintf("o.deliver_contact ILIKE $%d", argIndex),

			fmt.Sprintf("d.first_name ILIKE $%d", argIndex),
			fmt.Sprintf("d.last_name ILIKE $%d", argIndex),
			fmt.Sprintf("d.patronymic_name ILIKE $%d", argIndex),
			fmt.Sprintf("d.phone ILIKE $%d", argIndex),
			fmt.Sprintf("d.email ILIKE $%d", argIndex),
			fmt.Sprintf("d.meta ILIKE $%d", argIndex),

			fmt.Sprintf("v.numberplate ILIKE $%d", argIndex),
			fmt.Sprintf("v.trailer_numberplate ILIKE $%d", argIndex),
			fmt.Sprintf("v.meta ILIKE $%d", argIndex),
		}

		conditions = append(conditions, fmt.Sprintf("(%s)", strings.Join(searchConditions, " OR ")))
		args = append(args, searchTerm)
		argIndex++
	}

	if query.OfferCompanyID != nil && query.OfferExecCompanyID != nil {
		conditions = append(conditions, fmt.Sprintf("o.company_id = $%d", argIndex))
		args = append(args, *query.OfferCompanyID)
		argIndex++

		conditions = append(conditions, fmt.Sprintf("o.exec_company_id = $%d", argIndex))
		args = append(args, *query.OfferExecCompanyID)
		argIndex++
	} else if query.OfferCompanyID != nil {
		conditions = append(conditions, fmt.Sprintf("(o.company_id = $%d OR o.exec_company_id = $%d)", argIndex, argIndex))
		args = append(args, *query.OfferCompanyID)
		argIndex++
	} else if query.OfferExecCompanyID != nil {
		conditions = append(conditions, fmt.Sprintf("(o.company_id = $%d OR o.exec_company_id = $%d)", argIndex, argIndex))
		args = append(args, *query.OfferExecCompanyID)
		argIndex++
	}

	if query.OfferDriverID != nil {
		conditions = append(conditions, fmt.Sprintf("o.driver_id = $%d", argIndex))
		args = append(args, *query.OfferDriverID)
		argIndex++
	}

	if query.OfferVehicleID != nil {
		conditions = append(conditions, fmt.Sprintf("o.vehicle_id = $%d", argIndex))
		args = append(args, *query.OfferVehicleID)
		argIndex++
	}

	if query.OfferFromCountryID != nil {
		conditions = append(conditions, fmt.Sprintf("o.from_country_id = $%d", argIndex))
		args = append(args, *query.OfferFromCountryID)
		argIndex++
	}

	if query.OfferToCountryID != nil {
		conditions = append(conditions, fmt.Sprintf("o.to_country_id = $%d", argIndex))
		args = append(args, *query.OfferToCountryID)
		argIndex++
	}

	if query.OfferFromAddress != nil {
		conditions = append(conditions, fmt.Sprintf("o.from_address ILIKE $%d", argIndex))
		args = append(args, "%"+*query.OfferFromAddress+"%")
		argIndex++
	}

	if query.OfferToAddress != nil {
		conditions = append(conditions, fmt.Sprintf("o.to_address ILIKE $%d", argIndex))
		args = append(args, "%"+*query.OfferToAddress+"%")
		argIndex++
	}

	if query.OfferFromCountry != nil {
		conditions = append(conditions, fmt.Sprintf("o.from_country ILIKE $%d", argIndex))
		args = append(args, "%"+*query.OfferFromCountry+"%")
		argIndex++
	}

	if query.OfferToCountry != nil {
		conditions = append(conditions, fmt.Sprintf("o.to_country ILIKE $%d", argIndex))
		args = append(args, "%"+*query.OfferToCountry+"%")
		argIndex++
	}

	if query.OfferState != nil {
		conditions = append(conditions, fmt.Sprintf("o.offer_state = $%d", argIndex))
		args = append(args, *query.OfferState)
		argIndex++
	}

	if query.OfferRole != nil {
		conditions = append(conditions, fmt.Sprintf("o.offer_role = $%d", argIndex))
		args = append(args, *query.OfferRole)
		argIndex++
	}

	if query.OfferValidityStart != nil {
		conditions = append(conditions, fmt.Sprintf("o.validity_start >= $%d", argIndex))
		args = append(args, *query.OfferValidityStart)
		argIndex++
	}

	if query.OfferValidityEnd != nil {
		conditions = append(conditions, fmt.Sprintf("o.validity_end <= $%d", argIndex))
		args = append(args, *query.OfferValidityEnd)
		argIndex++
	}

	if query.OfferDeliveryStart != nil {
		conditions = append(conditions, fmt.Sprintf("o.delivery_start >= $%d", argIndex))
		args = append(args, *query.OfferDeliveryStart)
		argIndex++
	}

	if query.OfferDeliveryEnd != nil {
		conditions = append(conditions, fmt.Sprintf("o.delivery_end <= $%d", argIndex))
		args = append(args, *query.OfferDeliveryEnd)
		argIndex++
	}

	if query.OfferCostPerKmMin != nil {
		conditions = append(conditions, fmt.Sprintf("o.cost_per_km >= $%d", argIndex))
		args = append(args, *query.OfferCostPerKmMin)
		argIndex++
	}

	if query.OfferCostPerKmMax != nil {
		conditions = append(conditions, fmt.Sprintf("o.cost_per_km <= $%d", argIndex))
		args = append(args, *query.OfferCostPerKmMax)
		argIndex++
	}

	if query.OfferPriceMin != nil {
		conditions = append(conditions, fmt.Sprintf("o.offer_price >= $%d", argIndex))
		args = append(args, *query.OfferPriceMin)
		argIndex++
	}

	if query.OfferPriceMax != nil {
		conditions = append(conditions, fmt.Sprintf("o.offer_price <= $%d", argIndex))
		args = append(args, *query.OfferPriceMax)
		argIndex++
	}

	if query.OfferTotalPriceMin != nil {
		conditions = append(conditions, fmt.Sprintf("o.total_price >= $%d", argIndex))
		args = append(args, *query.OfferTotalPriceMin)
		argIndex++
	}

	if query.OfferTotalPriceMax != nil {
		conditions = append(conditions, fmt.Sprintf("o.total_price <= $%d", argIndex))
		args = append(args, *query.OfferTotalPriceMax)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	queryStr := fmt.Sprintf(`
        SELECT DISTINCT o.id 
        FROM tbl_offer o
        LEFT JOIN tbl_driver d ON o.driver_id = d.id AND d.active = 1 AND d.deleted = 0
        LEFT JOIN tbl_vehicle v ON o.vehicle_id = v.id AND v.deleted = 0
        %s`, whereClause)

	var offerIDs []int
	err := pgxscan.Select(context.Background(), db.DB, &offerIDs, queryStr, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get offer IDs: %w", err)
	}

	return offerIDs, nil
}

func hasOfferFilters(query dto.OfferFiltersQuery) bool {
	return query.OfferCompanyID != nil ||
		query.OfferExecCompanyID != nil ||
		query.OfferDriverID != nil ||
		query.OfferVehicleID != nil ||
		query.OfferFromCountryID != nil ||
		query.OfferToCountryID != nil ||
		query.OfferFromAddress != nil ||
		query.OfferToAddress != nil ||
		query.OfferFromCountry != nil ||
		query.OfferToCountry != nil ||
		query.OfferState != nil ||
		query.OfferRole != nil ||
		query.OfferValidityStart != nil ||
		query.OfferValidityEnd != nil ||
		query.OfferDeliveryStart != nil ||
		query.OfferDeliveryEnd != nil ||
		query.Search != nil ||
		query.OfferCostPerKmMin != nil ||
		query.OfferCostPerKmMax != nil ||
		query.OfferPriceMin != nil ||
		query.OfferPriceMax != nil ||
		query.OfferTotalPriceMin != nil ||
		query.OfferTotalPriceMax != nil
}
