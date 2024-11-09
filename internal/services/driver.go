package services

import (
	"context"
	"encoding/json"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	db "texApi/database"
	"texApi/internal/dto"
	"texApi/internal/queries"
	"texApi/pkg/utils"
)

func GetDriverList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
	offset := (page - 1) * perPage

	rows, err := db.DB.Query(
		context.Background(),
		queries.GetDriverList,
		perPage,
		offset,
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
			&driver.CreatedAt, &driver.UpdatedAt, &driver.Active,
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

	response := dto.PaginatedResponse{
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
	var companyJSON, vehiclesJSON []byte

	err := db.DB.QueryRow(
		context.Background(),
		queries.GetDriverByID,
		id,
	).Scan(
		&driver.ID, &driver.UUID, &driver.CompanyID, &driver.FirstName,
		&driver.LastName, &driver.PatronymicName, &driver.Phone,
		&driver.Email, &driver.Featured, &driver.Rating,
		&driver.Partner, &driver.SuccessfulOps, &driver.ImageURL,
		&driver.CreatedAt, &driver.UpdatedAt, &driver.Active,
		&companyJSON, &vehiclesJSON,
	)

	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Driver not found", err.Error()))
		return
	}

	json.Unmarshal(companyJSON, &driver.Company)
	json.Unmarshal(vehiclesJSON, &driver.AssignedVehicles)

	ctx.JSON(http.StatusOK, utils.FormatResponse("Driver details", driver))
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
		driver.CompanyID, driver.FirstName, driver.LastName,
		driver.PatronymicName, driver.Phone, driver.Email,
		driver.ImageURL,
	).Scan(&id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating driver", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.FormatResponse("Successfully created driver!", gin.H{"id": id}))
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
		id, driver.FirstName, driver.LastName, driver.PatronymicName,
		driver.Phone, driver.Email, driver.Featured, driver.Rating,
		driver.Partner, driver.ImageURL, driver.Active,
	).Scan(&updatedID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error updating driver", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully updated driver!", gin.H{"id": updatedID}))
}

func DeleteDriver(ctx *gin.Context) {
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

//
//func CreateDriver(ctx *gin.Context) {
//	var driver dto.DriverCreate
//
//	if err := ctx.ShouldBindJSON(&driver); err != nil {
//		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
//		return
//	}
//
//	var id int
//	err := db.DB.QueryRow(
//		context.Background(),
//		queries.CreateDriver,
//		driver.CompanyID,
//		driver.FirstName,
//		driver.LastName,
//		driver.PatronymicName,
//		driver.Phone,
//		driver.Email,
//		driver.AvatarURL,
//	).Scan(&id)
//
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating driver", err.Error()))
//		return
//	}
//
//	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully created!", gin.H{"id": id}))
//}
//
//func UpdateDriver(ctx *gin.Context) {
//	id := ctx.Param("id")
//	var driver dto.DriverUpdate
//
//	if err := ctx.ShouldBindJSON(&driver); err != nil {
//		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
//		return
//	}
//
//	var updatedID int
//	err := db.DB.QueryRow(
//		context.Background(),
//		queries.UpdateDriver,
//		id,
//		driver.FirstName,
//		driver.LastName,
//		driver.PatronymicName,
//		driver.Phone,
//		driver.Email,
//		driver.AvatarURL,
//	).Scan(&updatedID)
//
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error updating driver", err.Error()))
//		return
//	}
//
//	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully updated!", gin.H{"id": updatedID}))
//}
//
//func DeleteDriver(ctx *gin.Context) {
//	id := ctx.Param("id")
//
//	_, err := db.DB.Exec(
//		context.Background(),
//		queries.DeleteDriver,
//		id,
//	)
//
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting driver"})
//		return
//	}
//	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully deleted!", gin.H{"id": id}))
//}
