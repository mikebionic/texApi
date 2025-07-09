package services

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"texApi/internal/dto"
	"texApi/internal/repo"
	"texApi/pkg/utils"
)

func StartTrip(ctx *gin.Context) {
	var input dto.StartTripInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid input data", err.Error()))
		return
	}

	if len(input.Offers) == 0 {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("At least one offer is required", ""))
		return
	}

	tripID, err := repo.CreateTrip(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to start trip", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.FormatResponse("Trip started successfully", map[string]interface{}{
		"trip_id": tripID,
	}))
}

func EndTrip(ctx *gin.Context) {
	var input dto.EndTripInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid input data", err.Error()))
		return
	}

	err := repo.EndTrip(input)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "access denied") {
			ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Trip not found or access denied", err.Error()))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to end trip", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Trip ended successfully", nil))
}

func GetTrips(ctx *gin.Context) {
	var query dto.TripQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid query parameters", err.Error()))
		return
	}

	trips, err := repo.GetTrips(query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to retrieve trips", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Trips retrieved successfully", trips))
}

func CreateGPSLogs(ctx *gin.Context) {
	var logs []dto.GPSLogInput
	if err := ctx.ShouldBindJSON(&logs); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid input data", err.Error()))
		return
	}

	if len(logs) == 0 {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("At least one GPS log is required", ""))
		return
	}

	err := repo.CreateGPSLogs(logs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to create GPS logs", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.FormatResponse("GPS logs created successfully", map[string]interface{}{
		"count": len(logs),
	}))
}

func GetGPSLogs(ctx *gin.Context) {
	var query dto.GPSLogQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid query parameters", err.Error()))
		return
	}

	logs, err := repo.GetGPSLogs(query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to retrieve GPS logs", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("GPS logs retrieved successfully", logs))
}

func GetLastPositions(ctx *gin.Context) {
	var query dto.PositionQuery

	if tripIDs := ctx.QueryArray("trip_ids"); len(tripIDs) > 0 {
		query.TripIDs = parseIntArray(tripIDs)
	}
	if companyIDs := ctx.QueryArray("company_ids"); len(companyIDs) > 0 {
		query.CompanyIDs = parseIntArray(companyIDs)
	}
	if offerIDs := ctx.QueryArray("offer_ids"); len(offerIDs) > 0 {
		query.OfferIDs = parseIntArray(offerIDs)
	}
	if driverIDs := ctx.QueryArray("driver_ids"); len(driverIDs) > 0 {
		query.DriverIDs = parseIntArray(driverIDs)
	}
	if vehicleIDs := ctx.QueryArray("vehicle_ids"); len(vehicleIDs) > 0 {
		query.VehicleIDs = parseIntArray(vehicleIDs)
	}

	positions, err := repo.GetLastPositions(query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to retrieve last positions", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Last positions retrieved successfully", positions))
}

func parseIntArray(strArray []string) []int {
	var intArray []int
	for _, str := range strArray {
		if id, err := strconv.Atoi(str); err == nil {
			intArray = append(intArray, id)
		}
	}
	return intArray
}

func GetTripsDetailed(ctx *gin.Context) {
	var query dto.TripQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid query parameters", err.Error()))
		return
	}

	trips, err := repo.GetTripsDetailed(query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to retrieve detailed trips", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Detailed trips retrieved successfully", trips))
}
