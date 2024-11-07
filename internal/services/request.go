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

func CreateRequest(ctx *gin.Context) {
	var myRequest dto.RequestCreate
	if err := ctx.ShouldBindJSON(&myRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	// TODO: this also should be checked
	//myRequest.UserID = loginUser.ID
	//myRequest.CompanyID = loginUser.CompanyID

	//should it check wheter the driver and vehicle are from that company?

	var id int
	err := db.DB.QueryRow(
		context.Background(),
		queries.CreateMyRequest,
		&myRequest.UserID,
		&myRequest.CompanyID,
		&myRequest.DriverID,
		&myRequest.VehicleID,
		&myRequest.CostPerKM,
		&myRequest.FromCountry,
		&myRequest.FromRegion,
		&myRequest.ToCountry,
		&myRequest.ToRegion,
		&myRequest.ValidityStart,
		&myRequest.ValidityEnd,
		&myRequest.Note,
	).Scan(&id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating request", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully created!", gin.H{"id": id}))
}

func UpdateRequest(ctx *gin.Context) {
	userID, _ := strconv.Atoi(ctx.GetHeader("UserID"))
	id := ctx.Param("id")

	var myRequest dto.RequestUpdate
	if err := ctx.ShouldBindJSON(&myRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	// TODO: should driver and vehicle be checked if linked to company

	var updatedID int
	err := db.DB.QueryRow(
		context.Background(),
		queries.UpdateMyRequest,
		id,
		myRequest.DriverID,
		myRequest.VehicleID,
		myRequest.CostPerKM,
		myRequest.FromCountry,
		myRequest.FromRegion,
		myRequest.ToCountry,
		myRequest.ToRegion,
		myRequest.ValidityStart,
		myRequest.ValidityEnd,
		myRequest.Note,
		userID,
	).Scan(&updatedID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error updating request", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully updated!", gin.H{"id": updatedID}))
}

// User Company specific request
// TODO: take userID and validate in query
func GetCompanyRequests(ctx *gin.Context) {
	userID, _ := strconv.Atoi(ctx.GetHeader("UserID"))
	companyID, _ := strconv.Atoi(ctx.GetHeader("CompanyID"))

	stmt := queries.GetMyRequest + " AND company_id = $1 AND user_id = $2;"
	var myRequests []dto.RequestCreate

	err := pgxscan.Select(
		context.Background(), db.DB,
		&myRequests, stmt,
		companyID,
		userID,
	)

	if err != nil || len(myRequests) == 0 {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Requests not found", ""))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("My Requests", myRequests))
}

func GetRequests(ctx *gin.Context) {
	stmt := queries.GetMyRequest
	var allRequests []dto.RequestCreate

	err := pgxscan.Select(
		context.Background(), db.DB,
		&allRequests, stmt,
	)

	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Requests not found", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Requests", allRequests))
}

func DeleteRequest(ctx *gin.Context) {
	// TODO: validate user with request user
	userID, _ := strconv.Atoi(ctx.GetHeader("UserID"))
	id := ctx.Param("id")

	result, err := db.DB.Exec(
		context.Background(),
		queries.DeleteMyRequest,
		id,
		userID,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error deleting request", err.Error()))
		return
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Request not found or already deleted", "No matching request found"))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully deleted!", gin.H{"id": id}))
}
