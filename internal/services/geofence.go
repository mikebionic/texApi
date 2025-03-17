package services

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	db "texApi/database"
	"texApi/internal/dto"
	"texApi/pkg/utils"
)

func CreateGeofence(ctx *gin.Context) {
	var geofence dto.GeofenceCreate
	companyID := ctx.MustGet("companyID").(int)

	if err := ctx.ShouldBindJSON(&geofence); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	// Verify company ID from token matches request
	if geofence.CompanyID != companyID && ctx.MustGet("role").(string) != "admin" {
		ctx.JSON(http.StatusForbidden, utils.FormatErrorResponse("Permission denied", "You can only create geofences for your own company"))
		return
	}

	// Set default values
	isActive := true
	if geofence.IsActive != nil {
		isActive = *geofence.IsActive
	}

	// Prepare optional fields
	var description, radius interface{} = nil, nil
	if geofence.Description != nil {
		description = *geofence.Description
	}
	if geofence.Radius != nil {
		radius = *geofence.Radius
	}

	//// Validate geofence type and coordinates
	//if err := validateGeofence(geofence.FenceType, geofence.Coordinates); err != nil {
	//	ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid geofence", err.Error()))
	//	return
	//}

	// Insert geofence
	query := `
        INSERT INTO tbl_geofence (
            company_id, name, description, fence_type,
            coordinates, radius, is_active
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7
        ) RETURNING id, uuid
    `

	var geofenceID int
	var geofenceUUID string
	err := db.DB.QueryRow(
		context.Background(),
		query,
		geofence.CompanyID, geofence.Name, description, geofence.FenceType,
		geofence.Coordinates, radius, isActive,
	).Scan(&geofenceID, &geofenceUUID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Database error", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.FormatResponse("Geofence created successfully", gin.H{
		"geofence_id": geofenceID,
		"uuid":        geofenceUUID,
	}))
}

// GetGeofence retrieves details about a specific geofence
func GetGeofence(ctx *gin.Context) {
	geofenceID := ctx.Param("geofence_id")
	companyID := ctx.MustGet("companyID").(int)
	role := ctx.MustGet("role").(string)

	var geofence dto.GeofenceDetails
	var query string

	if role == "admin" {
		query = `
            SELECT * FROM tbl_geofence
            WHERE uuid = $1
        `
	} else {
		query = `
            SELECT * FROM tbl_geofence
            WHERE uuid = $1 AND company_id = $2
        `
	}

	var err error
	if role == "admin" {
		err = pgxscan.Get(context.Background(), db.DB, &geofence, query, geofenceID)
	} else {
		err = pgxscan.Get(context.Background(), db.DB, &geofence, query, geofenceID, companyID)
	}

	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Geofence not found", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Geofence details", geofence))
}

// GetGeofenceList retrieves a list of geofences
func GetGeofenceList(ctx *gin.Context) {
	companyID := ctx.MustGet("companyID").(int)
	role := ctx.MustGet("role").(string)
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
	offset := (page - 1) * perPage

	var countQuery, dataQuery string
	var args []interface{}

	if role == "admin" {
		// Optional company_id filter for admins
		companyIDFilter := ctx.Query("company_id")
		if companyIDFilter != "" {
			filterID, err := strconv.Atoi(companyIDFilter)
			if err == nil {
				countQuery = `
                    SELECT COUNT(*) FROM tbl_geofence
                    WHERE company_id = $1
                `
				dataQuery = `
                    SELECT *, COUNT(*) OVER() as total_count
                    FROM tbl_geofence
                    WHERE company_id = $1
                    ORDER BY name
                    LIMIT $2 OFFSET $3
                `
				args = []interface{}{filterID, perPage, offset}
			} else {
				ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid company ID", err.Error()))
				return
			}
		} else {
			countQuery = `SELECT COUNT(*) FROM tbl_geofence`
			dataQuery = `
                SELECT *, COUNT(*) OVER() as total_count
                FROM tbl_geofence
                ORDER BY name
                LIMIT $1 OFFSET $2
            `
			args = []interface{}{perPage, offset}
		}
	} else {
		// Regular users can only see their company's geofences
		countQuery = `
            SELECT COUNT(*) FROM tbl_geofence
            WHERE company_id = $1
        `
		dataQuery = `
            SELECT *, COUNT(*) OVER() as total_count
            FROM tbl_geofence
            WHERE company_id = $1
            ORDER BY name
            LIMIT $2 OFFSET $3
        `
		args = []interface{}{companyID, perPage, offset}
	}

	// Get count
	var totalCount int
	var countArgs []interface{}
	if role == "admin" && len(args) > 0 {
		countArgs = []interface{}{args[0]}
	} else if role != "admin" {
		countArgs = []interface{}{companyID}
	}

	if len(countArgs) > 0 {
		err := db.DB.QueryRow(context.Background(), countQuery, countArgs...).Scan(&totalCount)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Database error", err.Error()))
			return
		}
	} else {
		err := db.DB.QueryRow(context.Background(), countQuery).Scan(&totalCount)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Database error", err.Error()))
			return
		}
	}

	// Get data
	var geofences []dto.GeofenceDetails
	err := pgxscan.Select(context.Background(), db.DB, &geofences, dataQuery, args...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Database error", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Geofence list", gin.H{
		"total":     totalCount,
		"page":      page,
		"per_page":  perPage,
		"geofences": geofences,
		"pages":     (totalCount + perPage - 1) / perPage,
	}))
}

// UpdateGeofence updates an existing geofence
func UpdateGeofence(ctx *gin.Context) {
	geofenceID := ctx.Param("geofence_id")
	var geofenceUpdate dto.GeofenceUpdate
	companyID := ctx.MustGet("companyID").(int)
	role := ctx.MustGet("role").(string)

	if err := ctx.ShouldBindJSON(&geofenceUpdate); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	// Verify geofence belongs to the company
	var verifyQuery string
	if role == "admin" {
		verifyQuery = `
            SELECT COUNT(*) FROM tbl_geofence
            WHERE uuid = $1
        `
	} else {
		verifyQuery = `
            SELECT COUNT(*) FROM tbl_geofence
            WHERE uuid = $1 AND company_id = $2
        `
	}

	var count int
	var err error
	if role == "admin" {
		err = db.DB.QueryRow(context.Background(), verifyQuery, geofenceID).Scan(&count)
	} else {
		err = db.DB.QueryRow(context.Background(), verifyQuery, geofenceID, companyID).Scan(&count)
	}

	if err != nil || count == 0 {
		ctx.JSON(http.StatusForbidden, utils.FormatErrorResponse("Permission denied", "Geofence not found or does not belong to your company"))
		return
	}

	//// Validate geofence type and coordinates if provided
	//if geofenceUpdate.FenceType != nil && geofenceUpdate.Coordinates != nil {
	//	if err := validateGeofence(*geofenceUpdate.FenceType, geofenceUpdate.Coordinates); err != nil {
	//		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid geofence", err.Error()))
	//		return
	//	}
	//}

	// Build update query
	var setClauses []string
	var args []interface{}
	args = append(args, geofenceID) // $1 is the geofence UUID

	// Add fields to update
	if geofenceUpdate.Name != nil {
		setClauses = append(setClauses, "name = $"+strconv.Itoa(len(args)+1))
		args = append(args, *geofenceUpdate.Name)
	}

	if geofenceUpdate.Description != nil {
		setClauses = append(setClauses, "description = $"+strconv.Itoa(len(args)+1))
		args = append(args, *geofenceUpdate.Description)
	}

	if geofenceUpdate.FenceType != nil {
		setClauses = append(setClauses, "fence_type = $"+strconv.Itoa(len(args)+1))
		args = append(args, *geofenceUpdate.FenceType)
	}

	if geofenceUpdate.Coordinates != nil {
		setClauses = append(setClauses, "coordinates = $"+strconv.Itoa(len(args)+1))
		args = append(args, geofenceUpdate.Coordinates)
	}

	if geofenceUpdate.Radius != nil {
		setClauses = append(setClauses, "radius = $"+strconv.Itoa(len(args)+1))
		args = append(args, *geofenceUpdate.Radius)
	}

	if geofenceUpdate.IsActive != nil {
		setClauses = append(setClauses, "is_active = $"+strconv.Itoa(len(args)+1))
		args = append(args, *geofenceUpdate.IsActive)
	}

	// Always update updated_at
	setClauses = append(setClauses, "updated_at = CURRENT_TIMESTAMP")

	// If no fields to update, return
	if len(setClauses) == 1 {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("No fields to update", "Request body is empty"))
		return
	}

	// Execute update
	query := `
        UPDATE tbl_geofence
        SET ` + strings.Join(setClauses, ", ") + `
        WHERE uuid = $1
        RETURNING id
    `

	var updatedID int
	err = db.DB.QueryRow(context.Background(), query, args...).Scan(&updatedID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Database error", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Geofence updated successfully", gin.H{
		"geofence_id": updatedID,
		"uuid":        geofenceID,
	}))
}

// DeleteGeofence deletes a geofence
func DeleteGeofence(ctx *gin.Context) {
	geofenceID := ctx.Param("geofence_id")
	companyID := ctx.MustGet("companyID").(int)
	role := ctx.MustGet("role").(string)

	// Verify geofence belongs to the company
	var verifyQuery string
	if role == "admin" {
		verifyQuery = `
            SELECT id FROM tbl_geofence
            WHERE uuid = $1
        `
	} else {
		verifyQuery = `
            SELECT id FROM tbl_geofence
            WHERE uuid = $1 AND company_id = $2
        `
	}

	var id int
	var err error
	if role == "admin" {
		err = db.DB.QueryRow(context.Background(), verifyQuery, geofenceID).Scan(&id)
	} else {
		err = db.DB.QueryRow(context.Background(), verifyQuery, geofenceID, companyID).Scan(&id)
	}

	if err != nil {
		ctx.JSON(http.StatusForbidden, utils.FormatErrorResponse("Permission denied", "Geofence not found or does not belong to your company"))
		return
	}

	// Delete geofence
	query := `DELETE FROM tbl_geofence WHERE uuid = $1`
	_, err = db.DB.Exec(context.Background(), query, geofenceID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Database error", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Geofence deleted successfully", gin.H{
		"geofence_id": id,
		"uuid":        geofenceID,
	}))
}

//// GetGeofenceEvents retrieves geofence events
//func GetGeofenceEvents(ctx *gin.Context) {
//	geofenceIDStr := ctx.Query("geofence_id")
//	vehicleIDStr := ctx.Query("vehicle_id")
//	driverIDStr := ctx.Query("driver_id")
//	startTimeStr := ctx.Query("start_time")
//	endTimeStr := ctx.Query("end_time")
//	eventType := ctx.Query("event_type")
//	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
//	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
//	offset := (page - 1) * perPage
//
//	companyID := ctx.MustGet("companyID").(int)
//	role := ctx.MustGet("role").(string)
//
//	// Parse filter parameters
//	var geofenceID, vehicleID, driverID int
//	var startTime, endTime time.Time
//	var err error
//
//	if geofenceIDStr != "" {
//		geofenceID, err = strconv.Atoi(geofenceIDStr)
//		if err != nil {
//			ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid geofence ID", err.Error()))
//			return
