package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	db "texApi/database"
	"texApi/internal/dto"
	"texApi/internal/queries"
	"texApi/pkg/utils"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
)

func CreateVehicle(ctx *gin.Context) {
	var vehicle dto.VehicleCreate

	if err := ctx.ShouldBindJSON(&vehicle); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	companyID := ctx.MustGet("companyID").(int)
	role := ctx.MustGet("role")
	if !(role == "admin" || role == "system") {
		vehicle.CompanyID = companyID
	}

	var id int
	err := db.DB.QueryRow(
		context.Background(),
		queries.CreateVehicle,
		vehicle.CompanyID, vehicle.VehicleType, vehicle.VehicleBrandID,
		vehicle.VehicleModelID, vehicle.YearOfIssue, vehicle.Mileage,
		vehicle.Numberplate, vehicle.TrailerNumberplate, vehicle.Gps,
		vehicle.Photo1URL, vehicle.Photo2URL, vehicle.Photo3URL,
		vehicle.Docs1URL, vehicle.Docs2URL, vehicle.Docs3URL,
	).Scan(&id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating vehicle", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.FormatResponse("Successfully created vehicle!", gin.H{"id": id}))
}

func UpdateVehicle(ctx *gin.Context) {
	id := ctx.Param("id")
	var vehicle dto.VehicleUpdate

	stmt := queries.UpdateVehicle

	companyID := ctx.MustGet("companyID").(int)
	role := ctx.MustGet("role")
	if !(role == "admin" || role == "system") {
		vehicle.CompanyID = &companyID
		stmt += ` WHERE (id = $1 AND company_id = $17) AND (active = 1 AND deleted = 0)`
	} else {
		stmt += ` WHERE id = $1`
	}
	stmt += ` RETURNING id;`

	if err := ctx.ShouldBindJSON(&vehicle); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	var updatedID int
	err := db.DB.QueryRow(
		context.Background(),
		stmt,
		id, vehicle.VehicleType, vehicle.VehicleBrandID,
		vehicle.VehicleModelID, vehicle.YearOfIssue, vehicle.Mileage,
		vehicle.Numberplate, vehicle.TrailerNumberplate, vehicle.Gps,
		vehicle.Photo1URL, vehicle.Photo2URL, vehicle.Photo3URL,
		vehicle.Docs1URL, vehicle.Docs2URL, vehicle.Docs3URL,
		vehicle.Active, vehicle.CompanyID, vehicle.Deleted,
	).Scan(&updatedID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error updating vehicle", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully updated vehicle!", gin.H{"id": updatedID}))
}

func DeleteVehicle(ctx *gin.Context) {
	role := ctx.MustGet("role")
	if !(role == "admin" || role == "system") {
		ctx.JSON(http.StatusForbidden, utils.FormatErrorResponse("Operation can't be done by user", ""))
		return
	}

	id := ctx.Param("id")

	_, err := db.DB.Exec(
		context.Background(),
		queries.DeleteVehicle,
		id,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error deleting vehicle", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully deleted vehicle!", gin.H{"id": id}))
}
func GetVehicleList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
	offset := (page - 1) * perPage

	rows, err := db.DB.Query(
		context.Background(),
		queries.GetVehicleList,
		perPage,
		offset,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Database error", err.Error()))
		return
	}
	defer rows.Close()

	var vehicles []dto.VehicleDetails
	var totalCount int

	for rows.Next() {
		var vehicle dto.VehicleDetails
		var companyJSON, brandJSON, modelJSON []byte

		// Scan all the columns returned by the query
		err := rows.Scan(
			&vehicle.ID, &vehicle.UUID, &vehicle.CompanyID, &vehicle.VehicleType,
			&vehicle.VehicleBrandID, &vehicle.VehicleModelID, &vehicle.YearOfIssue,
			&vehicle.Mileage, &vehicle.Numberplate, &vehicle.TrailerNumberplate,
			&vehicle.Gps, &vehicle.Photo1URL, &vehicle.Photo2URL,
			&vehicle.Photo3URL, &vehicle.Docs1URL, &vehicle.Docs2URL,
			&vehicle.Docs3URL, &vehicle.ViewCount, &vehicle.CreatedAt,
			&vehicle.UpdatedAt, &vehicle.Active, &vehicle.Deleted, &totalCount,
			&companyJSON, &brandJSON, &modelJSON,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Scan error", err.Error()))
			return
		}

		// Unmarshal the JSON fields
		json.Unmarshal(companyJSON, &vehicle.Company)
		json.Unmarshal(brandJSON, &vehicle.Brand)
		json.Unmarshal(modelJSON, &vehicle.Model)
		vehicles = append(vehicles, vehicle)
	}

	// Prepare response
	response := utils.PaginatedResponse{
		Total:   totalCount,
		Page:    page,
		PerPage: perPage,
		Data:    vehicles,
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Vehicle list", response))
}
func GetVehicle(ctx *gin.Context) {
	id := ctx.Param("id")

	var vehicle dto.VehicleDetails
	var companyJSON, brandJSON, modelJSON []byte // These will hold the JSON data

	err := db.DB.QueryRow(
		context.Background(),
		queries.GetVehicleByID,
		id,
	).Scan(
		&vehicle.ID, &vehicle.UUID, &vehicle.CompanyID, &vehicle.VehicleType,
		&vehicle.VehicleBrandID, &vehicle.VehicleModelID, &vehicle.YearOfIssue,
		&vehicle.Mileage, &vehicle.Numberplate, &vehicle.TrailerNumberplate,
		&vehicle.Gps, &vehicle.Photo1URL, &vehicle.Photo2URL,
		&vehicle.Photo3URL, &vehicle.Docs1URL, &vehicle.Docs2URL,
		&vehicle.Docs3URL, &vehicle.ViewCount, &vehicle.Meta, &vehicle.Meta2, &vehicle.Meta3,
		&vehicle.CreatedAt, &vehicle.UpdatedAt, &vehicle.Active, &vehicle.Deleted,
		// These 3 JSON columns need to be scanned into the respective byte slices
		&companyJSON, &brandJSON, &modelJSON,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Vehicle not found", err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Database error", err.Error()))
		}
		return
	}

	// Unmarshal the JSON data into respective structs
	if err := json.Unmarshal(companyJSON, &vehicle.Company); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to parse company JSON", err.Error()))
		return
	}

	if err := json.Unmarshal(brandJSON, &vehicle.Brand); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to parse brand JSON", err.Error()))
		return
	}

	if err := json.Unmarshal(modelJSON, &vehicle.Model); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to parse model JSON", err.Error()))
		return
	}

	// Return the vehicle details
	ctx.JSON(http.StatusOK, utils.FormatResponse("Vehicle details", vehicle))
}

// Vehicle Brand Services
func SingleVehicleBrand(ctx *gin.Context) {
	id := ctx.Param("id")
	stmt := queries.GetVehicleBrand + " AND id = $1;"
	var brand []dto.VehicleBrand

	err := pgxscan.Select(
		context.Background(), db.DB,
		&brand, stmt, id,
	)

	if err != nil || len(brand) == 0 {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Vehicle brand not found", ""))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Vehicle brand", brand[0]))
}

func GetVehicleBrands(ctx *gin.Context) {
	var brands []dto.VehicleBrand

	err := pgxscan.Select(
		context.Background(), db.DB,
		&brands, queries.GetVehicleBrand,
	)

	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Vehicle brands not found", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Vehicle brands", brands))
}

func CreateVehicleBrand(ctx *gin.Context) {
	var brand dto.VehicleBrand

	if err := ctx.ShouldBindJSON(&brand); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	var id int
	err := db.DB.QueryRow(
		context.Background(),
		queries.CreateVehicleBrand,
		brand.Name,
		brand.Country,
		brand.FoundedYear,
	).Scan(&id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating vehicle brand", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully created!", gin.H{"id": id}))
}

func UpdateVehicleBrand(ctx *gin.Context) {
	id := ctx.Param("id")
	var brand dto.VehicleBrandUpdate

	if err := ctx.ShouldBindJSON(&brand); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	var updatedID int
	err := db.DB.QueryRow(
		context.Background(),
		queries.UpdateVehicleBrand,
		id,
		brand.Name,
		brand.Country,
		brand.FoundedYear,
	).Scan(&updatedID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error updating vehicle brand", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully updated!", gin.H{"id": updatedID}))
}

func DeleteVehicleBrand(ctx *gin.Context) {
	id := ctx.Param("id")

	_, err := db.DB.Exec(
		context.Background(),
		queries.DeleteVehicleBrand,
		id,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error deleting vehicle brand", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully deleted!", gin.H{"id": id}))
}

// Vehicle Type Services
func SingleVehicleType(ctx *gin.Context) {
	id := ctx.Param("id")
	stmt := queries.GetVehicleType + " AND id = $1;"
	var vehicleType []dto.VehicleType

	err := pgxscan.Select(
		context.Background(), db.DB,
		&vehicleType, stmt, id,
	)

	if err != nil || len(vehicleType) == 0 {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Vehicle type not found", ""))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Vehicle type", vehicleType[0]))
}

func GetVehicleTypes(ctx *gin.Context) {
	var types []dto.VehicleType

	err := pgxscan.Select(
		context.Background(), db.DB,
		&types, queries.GetVehicleType,
	)

	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Vehicle types not found", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Vehicle types", types))
}

func CreateVehicleType(ctx *gin.Context) {
	var vehicleType dto.VehicleType

	if err := ctx.ShouldBindJSON(&vehicleType); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	var id int
	err := db.DB.QueryRow(
		context.Background(),
		queries.CreateVehicleType,
		vehicleType.TitleEn,
		vehicleType.DescEn,
		vehicleType.TitleRu,
		vehicleType.DescRu,
		vehicleType.TitleTk,
		vehicleType.DescTk,
		vehicleType.TitleDe,
		vehicleType.DescDe,
		vehicleType.TitleAr,
		vehicleType.DescAr,
		vehicleType.TitleEs,
		vehicleType.DescEs,
		vehicleType.TitleFr,
		vehicleType.DescFr,
		vehicleType.TitleZh,
		vehicleType.DescZh,
		vehicleType.TitleJa,
		vehicleType.DescJa,
	).Scan(&id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating vehicle type", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully created!", gin.H{"id": id}))
}

func UpdateVehicleType(ctx *gin.Context) {
	id := ctx.Param("id")
	var vehicleType dto.VehicleTypeUpdate

	if err := ctx.ShouldBindJSON(&vehicleType); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	var updatedID int
	err := db.DB.QueryRow(
		context.Background(),
		queries.UpdateVehicleType,
		id,
		vehicleType.TitleEn,
		vehicleType.DescEn,
		vehicleType.TitleRu,
		vehicleType.DescRu,
		vehicleType.TitleTk,
		vehicleType.DescTk,
		vehicleType.TitleDe,
		vehicleType.DescDe,
		vehicleType.TitleAr,
		vehicleType.DescAr,
		vehicleType.TitleEs,
		vehicleType.DescEs,
		vehicleType.TitleFr,
		vehicleType.DescFr,
		vehicleType.TitleZh,
		vehicleType.DescZh,
		vehicleType.TitleJa,
		vehicleType.DescJa,
	).Scan(&updatedID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error updating vehicle type", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully updated!", gin.H{"id": updatedID}))
}

func DeleteVehicleType(ctx *gin.Context) {
	id := ctx.Param("id")

	_, err := db.DB.Exec(
		context.Background(),
		queries.DeleteVehicleType,
		id,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error deleting vehicle type", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully deleted!", gin.H{"id": id}))
}

// Vehicle Model Services
func SingleVehicleModel(ctx *gin.Context) {
	id := ctx.Param("id")
	stmt := queries.GetVehicleModel + " AND m.id = $1;"
	var model []dto.VehicleModel

	err := pgxscan.Select(
		context.Background(), db.DB,
		&model, stmt, id,
	)

	if err != nil || len(model) == 0 {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Vehicle model not found", ""))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Vehicle model", model[0]))
}

func GetVehicleModels(ctx *gin.Context) {
	var models []dto.VehicleModel

	err := pgxscan.Select(
		context.Background(), db.DB,
		&models, queries.GetVehicleModel,
	)

	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Vehicle models not found", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Vehicle models", models))
}

func CreateVehicleModel(ctx *gin.Context) {
	var model dto.VehicleModel

	if err := ctx.ShouldBindJSON(&model); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	var id int
	err := db.DB.QueryRow(
		context.Background(),
		queries.CreateVehicleModel,
		model.Name,
		model.Year,
		model.VehicleBrandID,
		model.VehicleTypeID,
		model.Feature,
	).Scan(&id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating vehicle model", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully created!", gin.H{"id": id}))
}

func UpdateVehicleModel(ctx *gin.Context) {
	id := ctx.Param("id")
	var model dto.VehicleModelUpdate

	if err := ctx.ShouldBindJSON(&model); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	var updatedID int
	err := db.DB.QueryRow(
		context.Background(),
		queries.UpdateVehicleModel,
		id,
		model.Name,
		model.Year,
		model.VehicleBrandID,
		model.VehicleTypeID,
		model.Feature,
	).Scan(&updatedID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error updating vehicle model", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully updated!", gin.H{"id": updatedID}))
}
func DeleteVehicleModel(ctx *gin.Context) {
	id := ctx.Param("id")

	_, err := db.DB.Exec(
		context.Background(),
		queries.DeleteVehicleModel,
		id,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error deleting vehicle model", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully deleted!", gin.H{"id": id}))
}
