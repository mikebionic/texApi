package services

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
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
		vehicle.CompanyID, vehicle.VehicleType, vehicle.Brand, vehicle.VehicleModel, vehicle.YearOfIssue, vehicle.Numberplate, vehicle.TrailerNumberplate, vehicle.GPSActive, vehicle.Photo1URL, vehicle.Photo2URL, vehicle.Photo3URL, vehicle.Docs1URL, vehicle.Docs2URL, vehicle.Docs3URL,
	).Scan(&id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating vehicle"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"vehicle_id": id})
}

func GetVehicle(ctx *gin.Context) {
	id := ctx.Param("id")

	var vehicle dto.VehicleCreate
	err := db.DB.QueryRow(
		context.Background(),
		queries.GetVehicle,
		id,
	).Scan(&vehicle.CompanyID, &vehicle.VehicleType, &vehicle.Brand, &vehicle.VehicleModel, &vehicle.YearOfIssue, &vehicle.Numberplate, &vehicle.TrailerNumberplate, &vehicle.GPSActive, &vehicle.Photo1URL, &vehicle.Photo2URL, &vehicle.Photo3URL, &vehicle.Docs1URL, &vehicle.Docs2URL, &vehicle.Docs3URL)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
		return
	}

	ctx.JSON(http.StatusOK, vehicle)
}

func GetVehicles(ctx *gin.Context) {
	companyID := ctx.Query("company_id")

	var vehicles []dto.VehicleGet
	rows, err := db.DB.Query(
		context.Background(),
		queries.GetVehicles,
		companyID,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving vehicles"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var vehicle dto.VehicleGet
		if err := rows.Scan(&vehicle.ID, &vehicle.CompanyID, &vehicle.VehicleType, &vehicle.Brand, &vehicle.VehicleModel, &vehicle.YearOfIssue, &vehicle.Numberplate, &vehicle.TrailerNumberplate, &vehicle.GPSActive, &vehicle.Photo1URL, &vehicle.Photo2URL, &vehicle.Photo3URL, &vehicle.Docs1URL, &vehicle.Docs2URL, &vehicle.Docs3URL, &vehicle.CreatedAt, &vehicle.UpdatedAt, &vehicle.Active, &vehicle.Deleted); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing vehicle data"})
			return
		}
		vehicles = append(vehicles, vehicle)
	}

	ctx.JSON(http.StatusOK, vehicles)
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating vehicle"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"vehicle_id": updatedID})
}

func DeleteVehicle(ctx *gin.Context) {
	id := ctx.Param("id")

	_, err := db.DB.Exec(
		context.Background(),
		queries.DeleteVehicle,
		id,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting vehicle"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Vehicle deleted"})
}
