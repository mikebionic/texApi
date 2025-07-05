package services

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"log"
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
						'company_id',v.company_id,
						'vehicle_type_id',v.vehicle_type_id,
						'vehicle_brand_id',v.vehicle_brand_id,
						'vehicle_model_id',v.vehicle_model_id,
						'year_of_issue',v.year_of_issue,
						'mileage',v.mileage,
						'numberplate',v.numberplate,
						'trailer_numberplate',v.trailer_numberplate,
						'gps',v.gps,
						'photo1_url',v.photo1_url,
						'photo2_url',v.photo2_url,
						'photo3_url',v.photo3_url,
						'docs1_url',v.docs1_url,
						'docs2_url',v.docs2_url,
						'docs3_url',v.docs3_url,
						'view_count',v.view_count,
						'meta',v.meta,
						'meta2',v.meta2,
						'meta3',v.meta3,
						'available',v.available
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
	companyID := ctx.MustGet("companyID").(int)

	var driver dto.DriverDetails

	err := pgxscan.Get(context.Background(), db.DB, &driver, queries.GetDriverByID, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Driver not found", err.Error()))
		return
	}

	if driver.CompanyID == companyID {
		var userCredentials dto.UserCredentials
		err = pgxscan.Get(context.Background(), db.DB, &userCredentials, queries.GetUserCredentialsByDriverID, id)
		if err != nil {
			log.Printf("Warning: Could not fetch user credentials for driver %s: %v", id, err)
		} else {
			driver.UserCredentials = &userCredentials
		}
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

	tx, err := db.DB.Begin(context.Background())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error starting transaction", err.Error()))
		return
	}
	defer tx.Rollback(context.Background())

	var driverID int
	err = pgxscan.Get(context.Background(), tx, &driverID, queries.CreateDriver,
		driver.CompanyID, driver.FirstName, driver.LastName,
		driver.PatronymicName, driver.Phone, driver.Email, driver.ImageURL,
		driver.Meta, driver.Meta2, driver.Meta3)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating driver", err.Error()))
		return
	}

	username := fmt.Sprintf("driver-%d", driverID)
	password := utils.GenerateOTP(8)

	var userID int
	err = pgxscan.Get(context.Background(), tx, &userID, queries.CreateDriverUser,
		username, password, driver.Email, driver.Phone, "driver", 6, 1, 1, 0, driverID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating driver user account", err.Error()))
		return
	}

	if err = tx.Commit(context.Background()); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error committing transaction", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.FormatResponse("Successfully created driver with user account!", gin.H{
		"user_id":   userID,
		"driver_id": driverID,
		"username":  username,
		"email":     driver.Email,
		"phone":     driver.Phone,
		"password":  password,
	}))
}

func UpdateDriver(ctx *gin.Context) {
	id := ctx.Param("id")
	driverID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid driver ID", err.Error()))
		return
	}

	var driver dto.DriverUpdate
	if err = ctx.ShouldBindJSON(&driver); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	companyID := ctx.MustGet("companyID").(int)
	role := ctx.MustGet("role")
	stmt := queries.UpdateDriver

	if !(role == "admin" || role == "system") {
		driver.CompanyID = &companyID
		stmt += ` WHERE (id = $1 AND company_id = $11) AND deleted = 0`
	} else {
		stmt += ` WHERE id = $1`
	}
	stmt += ` RETURNING id;`

	tx, err := db.DB.Begin(context.Background())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error starting transaction", err.Error()))
		return
	}
	defer tx.Rollback(context.Background())

	var updatedID int
	err = pgxscan.Get(context.Background(), tx, &updatedID, stmt, id,
		driver.FirstName, driver.LastName, driver.PatronymicName,
		driver.Phone, driver.Email, driver.ImageURL, driver.Meta,
		driver.Meta2, driver.Meta3, driver.CompanyID, driver.BlockReason,
		driver.Active, driver.Deleted)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error updating driver", err.Error()))
		return
	}

	_, err = tx.Exec(context.Background(), queries.UpdateDriverUser,
		driver.Email, driver.Phone, driver.Active, driver.Deleted, driverID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error updating driver user account", err.Error()))
		return
	}

	if err = tx.Commit(context.Background()); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error committing transaction", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully updated driver and user account!", gin.H{
		"id":    updatedID,
		"email": driver.Email,
		"phone": driver.Phone,
	}))
}

func DeleteDriver(ctx *gin.Context) {
	id := ctx.Param("id")
	driverID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid driver ID", err.Error()))
		return
	}

	stmt := queries.DeleteDriver
	role := ctx.MustGet("role").(string)
	if !(role == "admin" || role == "system") {
		companyID := ctx.MustGet("companyID").(int)
		stmt += fmt.Sprintf(` AND company_id = %d`, companyID)
	}

	// Start transaction
	tx, err := db.DB.Begin(context.Background())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error starting transaction", err.Error()))
		return
	}
	defer tx.Rollback(context.Background())

	// Delete driver (soft delete)
	result, err := tx.Exec(context.Background(), stmt, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error deleting driver", err.Error()))
		return
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Driver not found or no changes were made", ""))
		return
	}

	// Delete corresponding user account (soft delete)
	_, err = tx.Exec(context.Background(), queries.DeleteDriverUser, driverID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error deleting driver user account", err.Error()))
		return
	}

	// Commit transaction
	if err = tx.Commit(context.Background()); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error committing transaction", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully deleted driver and user account!", gin.H{"id": id}))
}

func GetFilteredDriverList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
	offset := (page - 1) * perPage

	companyID := ctx.DefaultQuery("company_id", "")
	name := ctx.DefaultQuery("name", "") // Combined search instead of separate first_name, last_name
	phone := ctx.DefaultQuery("phone", "")
	email := ctx.DefaultQuery("email", "")
	vehicleID := ctx.DefaultQuery("vehicle_id", "")
	numberplate := ctx.DefaultQuery("numberplate", "")
	vehicleTypeID := ctx.DefaultQuery("vehicle_type_id", "")
	vehicleModelID := ctx.DefaultQuery("vehicle_model_id", "")
	vehicleBrandID := ctx.DefaultQuery("vehicle_brand_id", "")
	available := ctx.DefaultQuery("available", "")
	active := ctx.DefaultQuery("active", "")
	createdAtFrom := ctx.DefaultQuery("created_at_from", "")
	createdAtTo := ctx.DefaultQuery("created_at_to", "")
	metaSearch := ctx.DefaultQuery("meta_search", "")
	orderBy := ctx.DefaultQuery("order_by", "date-new")

	var whereClauses []string
	var args []interface{}
	argIndex := 3

	if companyID != "" {
		whereClauses = append(whereClauses, "d.company_id = $"+strconv.Itoa(argIndex))
		args = append(args, companyID)
		argIndex++
	}

	if name != "" {
		whereClauses = append(whereClauses, "(d.first_name ILIKE $"+strconv.Itoa(argIndex)+" OR d.last_name ILIKE $"+strconv.Itoa(argIndex)+" OR d.patronymic_name ILIKE $"+strconv.Itoa(argIndex)+" OR CONCAT(d.first_name, ' ', d.last_name) ILIKE $"+strconv.Itoa(argIndex)+")")
		args = append(args, "%"+name+"%")
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

	if active != "" {
		whereClauses = append(whereClauses, "d.active = $"+strconv.Itoa(argIndex))
		args = append(args, active)
		argIndex++
	}

	// Vehicle-related filters using EXISTS for better performance
	if vehicleID != "" || numberplate != "" || vehicleTypeID != "" || vehicleModelID != "" || vehicleBrandID != "" {
		var vehicleFilters []string

		if vehicleID != "" {
			vehicleFilters = append(vehicleFilters, "v.id = $"+strconv.Itoa(argIndex))
			args = append(args, vehicleID)
			argIndex++
		}

		if numberplate != "" {
			vehicleFilters = append(vehicleFilters, "v.numberplate ILIKE $"+strconv.Itoa(argIndex))
			args = append(args, "%"+numberplate+"%")
			argIndex++
		}

		if vehicleTypeID != "" {
			vehicleFilters = append(vehicleFilters, "v.vehicle_type_id = $"+strconv.Itoa(argIndex))
			args = append(args, vehicleTypeID)
			argIndex++
		}

		if vehicleModelID != "" {
			vehicleFilters = append(vehicleFilters, "v.vehicle_model_id = $"+strconv.Itoa(argIndex))
			args = append(args, vehicleModelID)
			argIndex++
		}

		if vehicleBrandID != "" {
			vehicleFilters = append(vehicleFilters, "v.vehicle_brand_id = $"+strconv.Itoa(argIndex))
			args = append(args, vehicleBrandID)
			argIndex++
		}

		whereClauses = append(whereClauses, "EXISTS (SELECT 1 FROM tbl_vehicle v WHERE v.company_id = d.company_id AND v.deleted = 0 AND "+strings.Join(vehicleFilters, " AND ")+")")
	}

	if createdAtFrom != "" {
		whereClauses = append(whereClauses, "d.created_at >= $"+strconv.Itoa(argIndex))
		args = append(args, createdAtFrom+" 00:00:00")
		argIndex++
	}

	if createdAtTo != "" {
		whereClauses = append(whereClauses, "d.created_at <= $"+strconv.Itoa(argIndex))
		args = append(args, createdAtTo+" 23:59:59")
		argIndex++
	}

	if metaSearch != "" {
		whereClauses = append(whereClauses, "(d.meta ILIKE $"+strconv.Itoa(argIndex)+" OR d.meta2 ILIKE $"+strconv.Itoa(argIndex)+" OR d.meta3 ILIKE $"+strconv.Itoa(argIndex)+")")
		args = append(args, "%"+metaSearch+"%")
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
	case "name":
		orderByClause = "d.first_name ASC, d.last_name ASC"
	case "company":
		orderByClause = "c.company_name ASC"
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
						'company_id',v.company_id,
						'vehicle_type_id',v.vehicle_type_id,
						'vehicle_brand_id',v.vehicle_brand_id,
						'vehicle_model_id',v.vehicle_model_id,
						'year_of_issue',v.year_of_issue,
						'mileage',v.mileage,
						'numberplate',v.numberplate,
						'trailer_numberplate',v.trailer_numberplate,
						'gps',v.gps,
						'photo1_url',v.photo1_url,
						'photo2_url',v.photo2_url,
						'photo3_url',v.photo3_url,
						'docs1_url',v.docs1_url,
						'docs2_url',v.docs2_url,
						'docs3_url',v.docs3_url,
						'view_count',v.view_count,
						'meta',v.meta,
						'meta2',v.meta2,
						'meta3',v.meta3,
						'available',v.available
					)
				)
				FROM tbl_vehicle v
				WHERE v.company_id = d.company_id AND v.deleted = 0
			), '[]') as assigned_vehicles
		FROM tbl_driver d
		LEFT JOIN tbl_company c ON d.company_id = c.id
	`
	query += fmt.Sprintf(" WHERE %s ORDER BY %s LIMIT $1 OFFSET $2", whereClause, orderByClause)

	queryArgs := append([]interface{}{perPage, offset}, args...)

	var drivers []dto.DriverDetails
	err := pgxscan.Select(
		context.Background(),
		db.DB,
		&drivers,
		query,
		queryArgs...,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Database error", err.Error()))
		return
	}

	var totalCount int
	if len(drivers) > 0 {
		totalCount, _ = strconv.Atoi(drivers[0].TotalCount)
	}

	response := utils.PaginatedResponse{
		Total:   totalCount,
		Page:    page,
		PerPage: perPage,
		Data:    drivers,
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Filtered driver list", response))
}

func CheckDriverNotBlocked(ctx *gin.Context, id int) bool {
	var driver dto.DriverDetails
	err := pgxscan.Get(context.Background(), db.DB, &driver, queries.GetDriverByID, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Driver not found", err.Error()))
		return false
	}

	if driver.Active == 0 {
		ctx.JSON(http.StatusUnauthorized, utils.FormatErrorResponse("Driver is blocked", utils.SafeString(driver.BlockReason)))
		return false
	}
	return true
}
