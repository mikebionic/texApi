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
