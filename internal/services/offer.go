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

func CreateOffer(ctx *gin.Context) {
	var myOffer dto.OfferCreate
	if err := ctx.ShouldBindJSON(&myOffer); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	// TODO: this also should be checked
	//myOffer.UserID = loginUser.ID
	//myOffer.CompanyID = loginUser.CompanyID

	//should it check wheter the driver and vehicle are from that company?

	var id int
	err := db.DB.QueryRow(
		context.Background(),
		queries.CreateOffer,
		&myOffer.UserID,
		&myOffer.CompanyID,
		&myOffer.DriverID,
		&myOffer.VehicleID,
		&myOffer.CostPerKM,
		&myOffer.FromCountry,
		&myOffer.FromRegion,
		&myOffer.ToCountry,
		&myOffer.ToRegion,
		&myOffer.ValidityStart,
		&myOffer.ValidityEnd,
		&myOffer.Note,
	).Scan(&id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating request", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully created!", gin.H{"id": id}))
}

func UpdateOffer(ctx *gin.Context) {
	userID, _ := strconv.Atoi(ctx.GetHeader("UserID"))
	id := ctx.Param("id")

	var myOffer dto.OfferUpdate
	if err := ctx.ShouldBindJSON(&myOffer); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	// TODO: should driver and vehicle be checked if linked to company

	var updatedID int
	err := db.DB.QueryRow(
		context.Background(),
		queries.UpdateOffer,
		id,
		myOffer.DriverID,
		myOffer.VehicleID,
		myOffer.CostPerKM,
		myOffer.FromCountry,
		myOffer.FromRegion,
		myOffer.ToCountry,
		myOffer.ToRegion,
		myOffer.ValidityStart,
		myOffer.ValidityEnd,
		myOffer.Note,
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
func GetCompanyOffers(ctx *gin.Context) {
	userID, _ := strconv.Atoi(ctx.GetHeader("UserID"))
	companyID, _ := strconv.Atoi(ctx.GetHeader("CompanyID"))

	stmt := queries.GetOffer + " AND company_id = $1 AND user_id = $2;"
	var myOffers []dto.OfferCreate

	err := pgxscan.Select(
		context.Background(), db.DB,
		&myOffers, stmt,
		companyID,
		userID,
	)

	if err != nil || len(myOffers) == 0 {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Requests not found", ""))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("My Requests", myOffers))
}

func GetOffers(ctx *gin.Context) {
	stmt := queries.GetOffer
	var allRequests []dto.OfferCreate

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

func DeleteOffer(ctx *gin.Context) {
	// TODO: validate user with request user
	userID, _ := strconv.Atoi(ctx.GetHeader("UserID"))
	id := ctx.Param("id")

	result, err := db.DB.Exec(
		context.Background(),
		queries.DeleteOffer,
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
