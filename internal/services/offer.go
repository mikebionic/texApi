package services

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	db "texApi/database"
	"texApi/internal/dto"
	"texApi/internal/queries"
	"texApi/pkg/utils"
)

func GetOfferList(ctx *gin.Context) {
	companyID := ctx.MustGet("companyID").(int)
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
	offset := (page - 1) * perPage

	rows, err := db.DB.Query(
		context.Background(),
		queries.GetOfferList,
		companyID,
		perPage,
		offset,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Database error", err.Error()))
		return
	}
	defer rows.Close()

	var offers []dto.OfferDetails
	var totalCount int

	for rows.Next() {
		var offer dto.OfferDetails
		var companyJSON, driverJSON, vehicleJSON, cargoJSON []byte

		err := rows.Scan(
			&offer.ID, &offer.UUID, &offer.UserID, &offer.CompanyID, &offer.DriverID, &offer.VehicleID, &offer.CargoID,
			&offer.OfferState, &offer.OfferRole, &offer.CostPerKm, &offer.Currency, &offer.FromCountry, &offer.FromRegion, &offer.ToCountry, &offer.ToRegion,
			&offer.FromAddress, &offer.ToAddress, &offer.SenderContact, &offer.RecipientContact, &offer.DeliverContact,
			&offer.ViewCount, &offer.ValidityStart, &offer.ValidityEnd, &offer.DeliveryStart, &offer.DeliveryEnd,
			&offer.Note, &offer.Tax, &offer.Trade, &offer.PaymentMethod, &offer.Meta, &offer.Meta2, &offer.Meta3,
			&offer.Featured, &offer.Partner, &offer.CreatedAt, &offer.UpdatedAt, &offer.Active, &offer.Deleted,
			&totalCount, &companyJSON, &driverJSON, &vehicleJSON, &cargoJSON,
		)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Scan error", err.Error()))
			return
		}

		json.Unmarshal(companyJSON, &offer.Company)
		json.Unmarshal(driverJSON, &offer.AssignedDriver)
		json.Unmarshal(vehicleJSON, &offer.AssignedVehicle)
		//json.Unmarshal(cargoJSON, &offer.Cargo)

		offers = append(offers, offer)
	}

	if err := rows.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Database error", err.Error()))
		return
	}

	response := utils.PaginatedResponse{
		Total:   totalCount,
		Page:    page,
		PerPage: perPage,
		Data:    offers,
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Offer list", response))
}

func GetOffer(ctx *gin.Context) {
	id := ctx.Param("id")
	//companyID := ctx.MustGet("companyID").(int)

	var offer dto.OfferDetails
	var companyJSON, driverJSON, vehicleJSON, cargoJSON []byte

	err := db.DB.QueryRow(
		context.Background(),
		queries.GetOfferByID,
		id,
	).Scan(
		&offer.ID, &offer.UUID, &offer.UserID, &offer.CompanyID, &offer.DriverID, &offer.VehicleID, &offer.CargoID,
		&offer.OfferState, &offer.OfferRole, &offer.CostPerKm, &offer.Currency, &offer.FromCountry, &offer.FromRegion, &offer.ToCountry, &offer.ToRegion,
		&offer.FromAddress, &offer.ToAddress, &offer.SenderContact, &offer.RecipientContact, &offer.DeliverContact,
		&offer.ViewCount, &offer.ValidityStart, &offer.ValidityEnd, &offer.DeliveryStart, &offer.DeliveryEnd,
		&offer.Note, &offer.Tax, &offer.Trade, &offer.PaymentMethod, &offer.Meta, &offer.Meta2, &offer.Meta3,
		&offer.Featured, &offer.Partner, &offer.CreatedAt, &offer.UpdatedAt, &offer.Active, &offer.Deleted,
		&companyJSON, &driverJSON, &vehicleJSON, &cargoJSON,
	)

	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Offer not found", err.Error()))
		return
	}

	json.Unmarshal(companyJSON, &offer.Company)
	json.Unmarshal(driverJSON, &offer.AssignedDriver)
	json.Unmarshal(vehicleJSON, &offer.AssignedVehicle)
	//json.Unmarshal(cargoJSON, &offer.Cargo)

	ctx.JSON(http.StatusOK, utils.FormatResponse("Offer details", offer))
}

func CreateOffer(ctx *gin.Context) {
	var offer dto.OfferCreate
	if err := ctx.ShouldBindJSON(&offer); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	companyID := ctx.MustGet("companyID").(int)
	userID := ctx.MustGet("id").(int)
	role := ctx.MustGet("role").(string)
	offer.CompanyID = companyID
	offer.UserID = userID
	offer.OfferRole = role

	var id int
	err := db.DB.QueryRow(
		context.Background(),
		queries.CreateOffer,
		offer.UserID, offer.CompanyID, offer.DriverID, offer.VehicleID, offer.CargoID, offer.CostPerKm, offer.Currency,
		offer.FromCountry, offer.FromRegion, offer.ToCountry, offer.ToRegion, offer.FromAddress, offer.ToAddress,
		offer.SenderContact, offer.RecipientContact, offer.DeliverContact, offer.ValidityStart, offer.ValidityEnd,
		offer.DeliveryStart, offer.DeliveryEnd, offer.Note, offer.Tax, offer.Trade, offer.PaymentMethod, offer.Meta, offer.Meta2, offer.Meta3,
		offer.OfferRole,
	).Scan(&id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating offer", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.FormatResponse("Successfully created offer!", gin.H{"id": id}))
}

func UpdateOffer(ctx *gin.Context) {
	id := ctx.Param("id")
	var offer dto.OfferUpdate

	stmt := queries.UpdateOffer

	companyID := ctx.MustGet("companyID").(int)
	if err := ctx.ShouldBindJSON(&offer); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	var updatedID int
	err := db.DB.QueryRow(
		context.Background(),
		stmt,
		id, offer.DriverID, offer.VehicleID, offer.CargoID, offer.CostPerKm, offer.Currency,
		offer.FromCountry, offer.FromRegion, offer.ToCountry, offer.ToRegion, offer.FromAddress, offer.ToAddress,
		offer.SenderContact, offer.RecipientContact, offer.DeliverContact, offer.ValidityStart, offer.ValidityEnd,
		offer.DeliveryStart, offer.DeliveryEnd, offer.Note, offer.Tax, offer.Trade, offer.PaymentMethod,
		offer.Meta, offer.Meta2, offer.Meta3, offer.Active, offer.Deleted, companyID,
	).Scan(&updatedID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error updating offer", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully updated offer!", gin.H{"id": updatedID}))
}

func DeleteOffer(ctx *gin.Context) {
	id := ctx.Param("id")
	companyID := ctx.MustGet("companyID").(int)

	_, err := db.DB.Exec(
		context.Background(),
		queries.DeleteOffer,
		id, companyID,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error deleting offer", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully deleted offer!", gin.H{"id": id}))
}
