package services

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	db "texApi/database"
	"texApi/internal/dto"
	"texApi/internal/queries"
)

func CreateDriver(ctx *gin.Context) {
	var driver dto.DriverCreate

	if err := ctx.ShouldBindJSON(&driver); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var id int
	err := db.DB.QueryRow(
		context.Background(),
		queries.CreateDriver,
		driver.CompanyID, driver.FirstName, driver.LastName, driver.PatronymicName, driver.Phone, driver.Email, driver.AvatarURL,
	).Scan(&id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating driver"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"driver_id": id})
}

func GetDriver(ctx *gin.Context) {
	id := ctx.Param("id")

	var driver dto.DriverCreate
	err := db.DB.QueryRow(
		context.Background(),
		queries.GetDriver,
		id,
	).Scan(&driver.CompanyID, &driver.FirstName, &driver.LastName, &driver.PatronymicName, &driver.Phone, &driver.Email, &driver.AvatarURL)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
		return
	}

	ctx.JSON(http.StatusOK, driver)
}

func UpdateDriver(ctx *gin.Context) {
	id := ctx.Param("id")
	var driver dto.DriverUpdate

	if err := ctx.ShouldBindJSON(&driver); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updatedID int
	err := db.DB.QueryRow(
		context.Background(),
		queries.UpdateDriver,
		id, driver.FirstName, driver.LastName, driver.PatronymicName, driver.Phone, driver.Email, driver.AvatarURL,
	).Scan(&updatedID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating driver"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"driver_id": updatedID})
}

func DeleteDriver(ctx *gin.Context) {
	id := ctx.Param("id")

	_, err := db.DB.Exec(
		context.Background(),
		queries.DeleteDriver,
		id,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting driver"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Driver deleted"})
}
