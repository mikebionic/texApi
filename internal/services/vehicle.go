package services

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	db "texApi/database"
	"texApi/internal/dto"
	"texApi/internal/queries"
	"texApi/pkg/utils"
)

func CreateVehicle(ctx *gin.Context) {
	var vehicle dto.VehicleCreate

	if err := ctx.ShouldBindJSON(&vehicle); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	var id int
	err := db.DB.QueryRow(
		context.Background(),
		queries.CreateVehicle,
		&vehicle.CompanyID,
		&vehicle.VehicleType,
		&vehicle.Brand,
		&vehicle.VehicleModel,
		&vehicle.YearOfIssue,
		&vehicle.Numberplate,
		&vehicle.TrailerNumberplate,
		&vehicle.GPSActive,
		&vehicle.Photo1URL,
		&vehicle.Photo2URL,
		&vehicle.Photo3URL,
		&vehicle.Docs1URL,
		&vehicle.Docs2URL,
		&vehicle.Docs3URL,
	).Scan(&id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating vehicle", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully created!", gin.H{"id": id}))
}

func GetVehicle(ctx *gin.Context) {
	//// TODO: add company_id validation
	id := ctx.Param("id")
	stmt := queries.GetVehicle + " AND id = $1;"
	var vehicles []dto.VehicleCreate

	err := pgxscan.Select(
		context.Background(), db.DB,
		&vehicles, stmt,
		id,
	)
	if err != nil || len(vehicles) == 0 {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Not found", ""))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Company vehicles", vehicles))
	return
}

func GetVehicles(ctx *gin.Context) {
	//// TODO: Change this to user.company_id valid drivers, not headers
	companyID, _ := strconv.Atoi(ctx.GetHeader("CompanyID"))
	stmt := queries.GetVehicle + " AND (company_id = $1 OR $1 = 0);"
	var vehicles []dto.VehicleCreate

	err := pgxscan.Select(
		context.Background(), db.DB,
		&vehicles, stmt,
		companyID,
	)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Error retrieving vehicles", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Company vehicles", vehicles))
	return
}

func UpdateVehicle(ctx *gin.Context) {
	id := ctx.Param("id")
	var vehicle dto.VehicleUpdate

	if err := ctx.ShouldBindJSON(&vehicle); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	var updatedID int
	err := db.DB.QueryRow(
		context.Background(),
		queries.UpdateVehicle,
		id, vehicle.VehicleType, vehicle.Brand, vehicle.VehicleModel, vehicle.YearOfIssue, vehicle.Numberplate, vehicle.TrailerNumberplate, vehicle.GPSActive, vehicle.Photo1URL, vehicle.Photo2URL, vehicle.Photo3URL, vehicle.Docs1URL, vehicle.Docs2URL, vehicle.Docs3URL,
	).Scan(&updatedID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatResponse("Error updating vehicle", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully updated!", gin.H{"id": updatedID}))
}

func DeleteVehicle(ctx *gin.Context) {
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

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully deleted!", gin.H{"id": id}))
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
		vehicleType.TypeName,
		vehicleType.Description,
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
		vehicleType.TypeName,
		vehicleType.Description,
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
		model.Brand,
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
		model.Brand,
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
