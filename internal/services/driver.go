package services

import (
	"context"
	"encoding/json"
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

	var drivers []dto.DriverDetails

	query := queries.GetDriverList

	err := pgxscan.Select(context.Background(), db.DB, &drivers, query, perPage, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Database error", err.Error()))
		return
	}

	// Return the paginated response
	response := utils.PaginatedResponse{
		Total:   len(drivers), // Update this with your own logic for total count
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
	role := ctx.MustGet("role")
	if !(role == "admin" || role == "system") {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Operation can't be done by user", ""))
		return
	}

	id := ctx.Param("id")

	_, err := db.DB.Exec(
		context.Background(),
		queries.DeleteDriver,
		id,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error deleting driver", err.Error()))
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
	argIndex := 1

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
						'vehicle_type', v.vehicle_type,
						'numberplate', v.numberplate
					)
				)
				FROM tbl_vehicle v
				WHERE v.company_id = d.company_id AND v.deleted = 0
			), '[]') as assigned_vehicles
		FROM tbl_driver d
		LEFT JOIN tbl_company c ON d.company_id = c.id
		WHERE ` + whereClause + `
		ORDER BY ` + orderByClause + `
		LIMIT $1 OFFSET $2
	`

	rows, err := db.DB.Query(
		context.Background(),
		query,
		append([]interface{}{perPage, offset}, args...)...,
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
			&driver.Email, &driver.Featured, &driver.Rating,
			&driver.Partner, &driver.SuccessfulOps, &driver.ImageURL,
			&driver.Meta, &driver.Meta2, &driver.Meta3, &driver.CreatedAt,
			&driver.UpdatedAt, &driver.Active, &driver.Deleted,
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

//
//func CreateDriver(ctx *gin.Context) {
//	var driver dto.DriverCreate
//
//	if err := ctx.ShouldBindJSON(&driver); err != nil {
//		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
//		return
//	}
//
//	companyID := ctx.MustGet("companyID").(int)
//	role := ctx.MustGet("role")
//	if !(role == "admin" || role == "system") {
//		driver.CompanyID = companyID
//	}
//
//	var id int
//	err := db.DB.QueryRow(
//		context.Background(),
//		queries.CreateDriver,
//		driver.CompanyID, driver.FirstName, driver.LastName,
//		driver.PatronymicName, driver.Phone, driver.Email,
//		driver.ImageURL, driver.Meta, driver.Meta2, driver.Meta3,
//	).Scan(&id)
//
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating driver", err.Error()))
//		return
//	}
//
//	ctx.JSON(http.StatusCreated, utils.FormatResponse("Successfully created driver!", gin.H{"id": id}))
//}
//
//func UpdateDriver(ctx *gin.Context) {
//	id := ctx.Param("id")
//	var driver dto.DriverUpdate
//
//	stmt := queries.UpdateDriver
//
//	companyID := ctx.MustGet("companyID").(int)
//	role := ctx.MustGet("role")
//	if !(role == "admin" || role == "system") {
//		driver.CompanyID = &companyID
//		stmt += ` WHERE (id = $1 AND company_id = $11) AND (active = 1 AND deleted = 0)`
//	} else {
//		stmt += ` WHERE id = $1`
//	}
//	stmt += ` RETURNING id;`
//
//	if err := ctx.ShouldBindJSON(&driver); err != nil {
//		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
//		return
//	}
//
//	var updatedID int
//	err := db.DB.QueryRow(
//		context.Background(),
//		stmt,
//		id, driver.FirstName, driver.LastName, driver.PatronymicName,
//		driver.Phone, driver.Email, driver.ImageURL, driver.Meta, driver.Meta2, driver.Meta3,
//		driver.CompanyID, driver.Active, driver.Deleted,
//	).Scan(&updatedID)
//
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error updating driver", err.Error()))
//		return
//	}
//
//	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully updated driver!", gin.H{"id": updatedID}))
//}
//
//func DeleteDriver(ctx *gin.Context) {
//	role := ctx.MustGet("role")
//	if !(role == "admin" || role == "system") {
//		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Operation can't be done by user", ""))
//		return
//	}
//
//	id := ctx.Param("id")
//
//	_, err := db.DB.Exec(
//		context.Background(),
//		queries.DeleteDriver,
//		id,
//	)
//
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error deleting driver", err.Error()))
//		return
//	}
//
//	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully deleted driver!", gin.H{"id": id}))
//}
