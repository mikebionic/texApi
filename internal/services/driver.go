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

func SingleDriver(ctx *gin.Context) {
	id := ctx.Param("id")
	stmt := queries.GetDriver + " AND id = $1;"
	var driver []dto.DriverCreate

	err := pgxscan.Select(
		context.Background(), db.DB,
		&driver, stmt, id,
	)

	if err != nil || len(driver) == 0 {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Driver not found", ""))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Driver", driver[0]))
	return
}

func GetDrivers(ctx *gin.Context) {
	companyID, _ := strconv.Atoi(ctx.GetHeader("CompanyID"))
	stmt := queries.GetDriver + " AND (company_id = $1 OR $1 = 0);"
	var drivers []dto.DriverCreate

	err := pgxscan.Select(
		context.Background(), db.DB,
		&drivers, stmt, companyID,
	)

	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Driver not found", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Company drivers", drivers))
	return
}

func CreateDriver(ctx *gin.Context) {
	var driver dto.DriverCreate

	if err := ctx.ShouldBindJSON(&driver); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	var id int
	err := db.DB.QueryRow(
		context.Background(),
		queries.CreateDriver,
		driver.CompanyID,
		driver.FirstName,
		driver.LastName,
		driver.PatronymicName,
		driver.Phone,
		driver.Email,
		driver.AvatarURL,
	).Scan(&id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating driver", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully created!", gin.H{"id": id}))
}

func UpdateDriver(ctx *gin.Context) {
	id := ctx.Param("id")
	var driver dto.DriverUpdate

	if err := ctx.ShouldBindJSON(&driver); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	var updatedID int
	err := db.DB.QueryRow(
		context.Background(),
		queries.UpdateDriver,
		id,
		driver.FirstName,
		driver.LastName,
		driver.PatronymicName,
		driver.Phone,
		driver.Email,
		driver.AvatarURL,
	).Scan(&updatedID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error updating driver", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully updated!", gin.H{"id": updatedID}))
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
	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully deleted!", gin.H{"id": id}))
}
