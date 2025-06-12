package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	db "texApi/database"
	"texApi/internal/dto"
	"texApi/internal/queries"
	"texApi/pkg/utils"
)

func GetDriverList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
	offset := (page - 1) * perPage

	search := ctx.Query("search")
	companyID := ctx.Query("company_id")
	featured := ctx.Query("featured")
	partner := ctx.Query("partner")
	rating := ctx.Query("min_rating")
	available := ctx.Query("available")
	active := ctx.Query("active")
	orderBy := ctx.DefaultQuery("order_by", "id")
	orderDir := ctx.DefaultQuery("order_dir", "ASC")
	includeVehicles := ctx.DefaultQuery("include_vehicles", "1")

	query := `
	WITH driver_data AS (
		SELECT
			d.*,
			COUNT(*) OVER() as total_count
		FROM tbl_driver d
		WHERE d.deleted = 0
	`

	args := make([]interface{}, 0)
	paramCount := 0
	if search != "" {
		query += fmt.Sprintf(`
			AND (LOWER(d.first_name) LIKE LOWER($%d)
			OR LOWER(d.last_name) LIKE LOWER($%d)
			OR LOWER(d.patronymic_name) LIKE LOWER($%d)
			OR LOWER(d.phone) LIKE LOWER($%d)
			OR LOWER(d.email) LIKE LOWER($%d)
			OR LOWER(d.meta) LIKE LOWER($%d)
			OR LOWER(d.meta2) LIKE LOWER($%d)
			OR LOWER(d.meta3) LIKE LOWER($%d))
		`, paramCount+1, paramCount+1, paramCount+1, paramCount+1, paramCount+1, paramCount+1, paramCount+1, paramCount+1)
		args = append(args, "%"+search+"%")
		paramCount++
	}

	if companyID != "" {
		query += fmt.Sprintf(" AND d.company_id = $%d", paramCount+1)
		args = append(args, companyID)
		paramCount++
	}

	if featured != "" {
		query += fmt.Sprintf(" AND d.featured = $%d", paramCount+1)
		args = append(args, featured)
		paramCount++
	}

	if partner != "" {
		query += fmt.Sprintf(" AND d.partner = $%d", paramCount+1)
		args = append(args, partner)
		paramCount++
	}

	if rating != "" {
		query += fmt.Sprintf(" AND d.rating >= $%d", paramCount+1)
		args = append(args, rating)
		paramCount++
	}

	if available != "" {
		query += fmt.Sprintf(" AND d.available = $%d", paramCount+1)
		args = append(args, available)
		paramCount++
	}

	if active != "" {
		query += fmt.Sprintf(" AND d.active = $%d", paramCount+1)
		args = append(args, active)
		paramCount++
	}

	validOrderColumns := map[string]bool{
		"id": true, "first_name": true, "last_name": true, "rating": true,
		"view_count": true, "successful_ops": true, "created_at": true, "updated_at": true,
	}

	if validOrderColumns[strings.ToLower(orderBy)] {
		query += fmt.Sprintf(" ORDER BY d.%s %s", orderBy, orderDir)
	} else {
		query += " ORDER BY d.id ASC"
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramCount+1, paramCount+2)
	args = append(args, perPage, offset)
	paramCount += 2

	query += `
	)
	SELECT
		dd.id, dd.uuid, dd.company_id, dd.first_name, dd.last_name,
		dd.patronymic_name, dd.phone, dd.email, dd.featured, dd.rating,
		dd.partner, dd.successful_ops, dd.image_url, dd.meta, dd.meta2,
		dd.meta3, dd.available, dd.view_count, dd.created_at, dd.updated_at,
		dd.active, dd.deleted, dd.total_count,
		json_build_object(
			'id', c.id,
			'company_name', c.company_name,
			'country', c.country
		) as company
	`

	if includeVehicles == "1" {
		query += `,
		COALESCE(
			(
				SELECT json_agg(
					json_build_object(
						'id', v.id,
						'vehicle_type_id', v.vehicle_type_id,
						'vehicle_brand_id', v.vehicle_brand_id,
						'numberplate', v.numberplate
					)
				)
				FROM tbl_vehicle v
				WHERE v.company_id = dd.company_id AND v.deleted = 0
			),
			'[]'
		) as assigned_vehicles
		`
	} else {
		query += ", '[]'::json as assigned_vehicles"
	}

	query += `
	FROM driver_data dd
	LEFT JOIN tbl_company c ON dd.company_id = c.id
	`

	var drivers []dto.DriverDetails
	err := pgxscan.Select(
		context.Background(),
		db.DB,
		&drivers,
		query,
		args...,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Database error", err.Error()))
		return
	}

	totalCount := 0
	if len(drivers) > 0 {
		totalCount, _ = strconv.Atoi(drivers[0].TotalCount)
	}

	response := utils.PaginatedResponse{
		Total:   totalCount,
		Page:    page,
		PerPage: perPage,
		Data:    drivers,
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Driver list", response))
}

func GetDriver(ctx *gin.Context) {
	id := ctx.Param("id")

	var driver dto.DriverDetails

	err := pgxscan.Get(context.Background(), db.DB, &driver, queries.GetDriverByID, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Driver not found", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Driver details", driver))
}

func CreateDriver(ctx *gin.Context) {
	var driver dto.DriverCreate
	if err := ctx.ShouldBindJSON(&driver); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	companyID := ctx.MustGet("companyID").(int)
	role := ctx.MustGet("role")
	if !(role == "admin" || role == "system") {
		driver.CompanyID = companyID
	}

	var id int
	err := pgxscan.Get(context.Background(), db.DB, &id, queries.CreateDriver, driver.CompanyID, driver.FirstName, driver.LastName,
		driver.PatronymicName, driver.Phone, driver.Email, driver.ImageURL, driver.Meta, driver.Meta2, driver.Meta3)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating driver", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.FormatResponse("Successfully created driver!", gin.H{"id": id}))
}

func UpdateDriver(ctx *gin.Context) {
	id := ctx.Param("id")
	var driver dto.DriverUpdate

	companyID := ctx.MustGet("companyID").(int)
	role := ctx.MustGet("role")
	stmt := queries.UpdateDriver

	if !(role == "admin" || role == "system") {
		driver.CompanyID = &companyID
		stmt += ` WHERE (id = $1 AND company_id = $11) AND (active = 1 AND deleted = 0)`
	} else {
		stmt += ` WHERE id = $1`
	}
	stmt += ` RETURNING id;`

	if err := ctx.ShouldBindJSON(&driver); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	var updatedID int
	err := pgxscan.Get(context.Background(), db.DB, &updatedID, stmt, id, driver.FirstName, driver.LastName, driver.PatronymicName,
		driver.Phone, driver.Email, driver.ImageURL, driver.Meta, driver.Meta2, driver.Meta3, driver.CompanyID, driver.Active, driver.Deleted)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error updating driver", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully updated driver!", gin.H{"id": updatedID}))
}

func DeleteDriver(ctx *gin.Context) {
	id := ctx.Param("id")
	stmt := queries.DeleteDriver

	role := ctx.MustGet("role").(string)
	if !(role == "admin" || role == "system") {
		companyID := ctx.MustGet("companyID").(int)
		stmt += fmt.Sprintf(` AND company_id = %d`, companyID)
	}

	result, err := db.DB.Exec(
		context.Background(),
		stmt,
		id,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error deleting driver", err.Error()))
		return
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Driver not found or no changes were made", ""))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully deleted driver!", gin.H{"id": id}))
}

func GetFilteredDriverList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
	offset := (page - 1) * perPage

	companyID := ctx.DefaultQuery("company_id", "")
	firstName := ctx.DefaultQuery("first_name", "")
	lastName := ctx.DefaultQuery("last_name", "")
	phone := ctx.DefaultQuery("phone", "")
	email := ctx.DefaultQuery("email", "")
	available := ctx.DefaultQuery("available", "")
	orderBy := ctx.DefaultQuery("order_by", "date-new")

	var whereClauses []string
	var args []interface{}
	argIndex := 3 // Start from 3 since $1 and $2 are reserved for perPage and offset

	if companyID != "" {
		whereClauses = append(whereClauses, "d.company_id = $"+strconv.Itoa(argIndex))
		args = append(args, companyID)
		argIndex++
	}
	if firstName != "" {
		whereClauses = append(whereClauses, "d.first_name ILIKE $"+strconv.Itoa(argIndex))
		args = append(args, "%"+firstName+"%")
		argIndex++
	}
	if lastName != "" {
		whereClauses = append(whereClauses, "d.last_name ILIKE $"+strconv.Itoa(argIndex))
		args = append(args, "%"+lastName+"%")
		argIndex++
	}
	if phone != "" {
		whereClauses = append(whereClauses, "d.phone ILIKE $"+strconv.Itoa(argIndex))
		args = append(args, "%"+phone+"%")
		argIndex++
	}
	if email != "" {
		whereClauses = append(whereClauses, "d.email ILIKE $"+strconv.Itoa(argIndex))
		args = append(args, "%"+email+"%")
		argIndex++
	}
	if available != "" {
		whereClauses = append(whereClauses, "d.available = $"+strconv.Itoa(argIndex))
		args = append(args, available)
		argIndex++
	}

	whereClause := "d.deleted = 0"
	if len(whereClauses) > 0 {
		whereClause += " AND " + strings.Join(whereClauses, " AND ")
	}

	var orderByClause string
	switch orderBy {
	case "date-old":
		orderByClause = "d.created_at ASC"
	case "date-new":
		orderByClause = "d.created_at DESC"
	case "id":
		orderByClause = "d.id ASC"
	case "rating":
		orderByClause = "d.rating DESC"
	case "successful_ops":
		orderByClause = "d.successful_ops DESC"
	case "view_count":
		orderByClause = "d.view_count DESC"
	default:
		orderByClause = "d.created_at DESC"
	}

	query := `
		SELECT 
			d.*, 
			COUNT(*) OVER() as total_count,
			json_build_object(
				'id', c.id,
				'company_name', c.company_name,
				'country', c.country
			) as company,
			COALESCE((
				SELECT json_agg(
					json_build_object(
						'id', v.id,
						'vehicle_type_id', v.vehicle_type_id,
						'vehicle_brand_id', v.vehicle_brand_id,
						'numberplate', v.numberplate
					)
				)
				FROM tbl_vehicle v
				WHERE v.company_id = d.company_id AND v.deleted = 0
			), '[]') as assigned_vehicles
		FROM tbl_driver d
		LEFT JOIN tbl_company c ON d.company_id = c.id
	`
	query += fmt.Sprintf(" WHERE %s ORDER BY %s LIMIT $1 OFFSET $2", whereClause, orderByClause)

	// First add perPage and offset to the args slice
	queryArgs := append([]interface{}{perPage, offset}, args...)

	rows, err := db.DB.Query(
		context.Background(),
		query,
		queryArgs...,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Database error", err.Error()))
		return
	}
	defer rows.Close()

	var drivers []dto.DriverDetails
	var totalCount int

	for rows.Next() {
		var driver dto.DriverDetails
		var companyJSON, vehiclesJSON []byte

		err := rows.Scan(
			&driver.ID, &driver.UUID, &driver.CompanyID, &driver.FirstName,
			&driver.LastName, &driver.PatronymicName, &driver.Phone,
			&driver.Email, &driver.Featured, &driver.Rating, &driver.Partner,
			&driver.SuccessfulOps, &driver.ImageURL, &driver.Meta, &driver.Meta2,
			&driver.Meta3, &driver.Available, &driver.ViewCount,
			&driver.CreatedAt, &driver.UpdatedAt, &driver.Active, &driver.Deleted,
			&totalCount, &companyJSON, &vehiclesJSON,
		)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Scan error", err.Error()))
			return
		}

		json.Unmarshal(companyJSON, &driver.Company)
		json.Unmarshal(vehiclesJSON, &driver.AssignedVehicles)
		drivers = append(drivers, driver)
	}

	// Prepare the response
	response := utils.PaginatedResponse{
		Total:   totalCount,
		Page:    page,
		PerPage: perPage,
		Data:    drivers,
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Filtered driver list", response))
}
