package services

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"net/http"
	db "texApi/database"
	"texApi/internal/dto"
	"texApi/pkg/utils"
)

func GetPriceQuoteList(ctx *gin.Context) {
	var filters dto.PriceQuoteFilters
	if err := ctx.ShouldBindQuery(&filters); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid query parameters", err.Error()))
		return
	}

	// Set defaults
	if filters.Page <= 0 {
		filters.Page = 1
	}
	if filters.PerPage <= 0 {
		filters.PerPage = 20
	}
	if filters.PerPage > 100 {
		filters.PerPage = 100
	}
	if filters.SortBy == "" {
		filters.SortBy = "created_at"
	}
	if filters.SortOrder == "" {
		filters.SortOrder = "DESC"
	}

	quotes, total, err := GetPriceQuotes(filters)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to retrieve price quotes", err.Error()))
		return
	}

	response := map[string]interface{}{
		"data":        quotes,
		"total":       total,
		"page":        filters.Page,
		"per_page":    filters.PerPage,
		"total_pages": (total + filters.PerPage - 1) / filters.PerPage,
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Price quotes retrieved successfully", response))
}

func CreatePriceQuote(ctx *gin.Context) {
	var req dto.CreatePriceQuoteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request data", err.Error()))
		return
	}

	quote, err := CreatePriceQuoteRecord(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to create price quote", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.FormatResponse("Price quote created successfully", quote))
}

func UpdatePriceQuote(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid price quote ID", err.Error()))
		return
	}

	var req dto.UpdatePriceQuoteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request data", err.Error()))
		return
	}

	quote, err := UpdatePriceQuoteRecord(id, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to update price quote", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Price quote updated successfully", quote))
}

func DeletePriceQuote(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid price quote ID", err.Error()))
		return
	}

	err = DeletePriceQuoteRecord(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to delete price quote", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Price quote deleted successfully", nil))
}

func GetPriceQuotes(filters dto.PriceQuoteFilters) ([]dto.PriceQuote, int, error) {
	var quotes []dto.PriceQuote
	var total int

	whereParts := []string{"deleted = 0"}
	args := []interface{}{}
	argIndex := 1

	if filters.TransportType != "" {
		whereParts = append(whereParts, fmt.Sprintf("transport_type = $%d", argIndex))
		args = append(args, filters.TransportType)
		argIndex++
	}
	if filters.SubType != "" {
		whereParts = append(whereParts, fmt.Sprintf("sub_type = $%d", argIndex))
		args = append(args, filters.SubType)
		argIndex++
	}
	if filters.FromCountry != "" {
		whereParts = append(whereParts, fmt.Sprintf("from_country ILIKE $%d", argIndex))
		args = append(args, "%"+filters.FromCountry+"%")
		argIndex++
	}
	if filters.ToCountry != "" {
		whereParts = append(whereParts, fmt.Sprintf("to_country ILIKE $%d", argIndex))
		args = append(args, "%"+filters.ToCountry+"%")
		argIndex++
	}
	if filters.FromRegion != "" {
		whereParts = append(whereParts, fmt.Sprintf("from_region ILIKE $%d", argIndex))
		args = append(args, "%"+filters.FromRegion+"%")
		argIndex++
	}
	if filters.ToRegion != "" {
		whereParts = append(whereParts, fmt.Sprintf("to_region ILIKE $%d", argIndex))
		args = append(args, "%"+filters.ToRegion+"%")
		argIndex++
	}
	if filters.Currency != "" {
		whereParts = append(whereParts, fmt.Sprintf("currency = $%d", argIndex))
		args = append(args, filters.Currency)
		argIndex++
	}
	if filters.PriceUnit != "" {
		whereParts = append(whereParts, fmt.Sprintf("price_unit = $%d", argIndex))
		args = append(args, filters.PriceUnit)
		argIndex++
	}
	if filters.MinPrice > 0 {
		whereParts = append(whereParts, fmt.Sprintf("average_price >= $%d", argIndex))
		args = append(args, filters.MinPrice)
		argIndex++
	}
	if filters.MaxPrice > 0 {
		whereParts = append(whereParts, fmt.Sprintf("average_price <= $%d", argIndex))
		args = append(args, filters.MaxPrice)
		argIndex++
	}
	if filters.FuelIncluded != nil {
		whereParts = append(whereParts, fmt.Sprintf("fuel_included = $%d", argIndex))
		args = append(args, *filters.FuelIncluded)
		argIndex++
	}
	if filters.CustomsIncluded != nil {
		whereParts = append(whereParts, fmt.Sprintf("customs_included = $%d", argIndex))
		args = append(args, *filters.CustomsIncluded)
		argIndex++
	}
	if filters.IsPromotional != nil {
		whereParts = append(whereParts, fmt.Sprintf("is_promotional = $%d", argIndex))
		args = append(args, *filters.IsPromotional)
		argIndex++
	}
	if filters.IsDynamic != nil {
		whereParts = append(whereParts, fmt.Sprintf("is_dynamic = $%d", argIndex))
		args = append(args, *filters.IsDynamic)
		argIndex++
	}
	if filters.DataSource != "" {
		whereParts = append(whereParts, fmt.Sprintf("data_source = $%d", argIndex))
		args = append(args, filters.DataSource)
		argIndex++
	}
	if filters.UserID > 0 {
		whereParts = append(whereParts, fmt.Sprintf("user_id = $%d", argIndex))
		args = append(args, filters.UserID)
		argIndex++
	}
	if filters.CompanyID > 0 {
		whereParts = append(whereParts, fmt.Sprintf("company_id = $%d", argIndex))
		args = append(args, filters.CompanyID)
		argIndex++
	}
	if filters.VehicleTypeID > 0 {
		whereParts = append(whereParts, fmt.Sprintf("vehicle_type_id = $%d", argIndex))
		args = append(args, filters.VehicleTypeID)
		argIndex++
	}
	if filters.Active != nil {
		whereParts = append(whereParts, fmt.Sprintf("active = $%d", argIndex))
		args = append(args, *filters.Active)
		argIndex++
	}

	whereClause := strings.Join(whereParts, " AND ")

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM tbl_price_quote WHERE %s", whereClause)
	err := pgxscan.Get(context.Background(), db.DB, &total, countQuery, args...)
	if err != nil {
		return quotes, 0, err
	}

	// Validate sort column
	validSortColumns := map[string]bool{
		"id": true, "created_at": true, "updated_at": true, "average_price": true,
		"min_price": true, "max_price": true, "transport_type": true,
		"from_country": true, "to_country": true, "validity_start": true, "validity_end": true,
	}
	if !validSortColumns[filters.SortBy] {
		filters.SortBy = "created_at"
	}

	if filters.SortOrder != "ASC" && filters.SortOrder != "DESC" {
		filters.SortOrder = "DESC"
	}

	// Get paginated results
	offset := (filters.Page - 1) * filters.PerPage
	query := fmt.Sprintf(`
		SELECT * FROM tbl_price_quote 
		WHERE %s 
		ORDER BY %s %s 
		LIMIT $%d OFFSET $%d`,
		whereClause, filters.SortBy, filters.SortOrder, argIndex, argIndex+1)

	args = append(args, filters.PerPage, offset)
	err = pgxscan.Select(context.Background(), db.DB, &quotes, query, args...)

	return quotes, total, err
}

func CreatePriceQuoteRecord(req dto.CreatePriceQuoteRequest) (dto.PriceQuote, error) {
	var quote dto.PriceQuote

	// Set defaults
	if req.ValidityStart.IsZero() {
		req.ValidityStart = time.Now()
	}
	if req.ValidityEnd.IsZero() {
		req.ValidityEnd = req.ValidityStart.AddDate(1, 0, 0) // +1 year
	}
	if req.Currency == "" {
		req.Currency = "USD"
	}
	if req.DataSource == "" {
		req.DataSource = "manual"
	}

	query := `
		INSERT INTO tbl_price_quote (
			transport_type, sub_type, user_id, company_id, exec_company_id, vehicle_type_id,
			packaging_type_id, cost_per_km, currency, from_country_id, from_city_id,
			to_country_id, to_city_id, distance, from_country, from_region, to_country, to_region,
			from_address, to_address, tax, tax_price, trade, discount, payment_method, payment_term,
			distance_km, average_price, min_price, max_price, price_unit, min_volume, max_volume,
			validity_start, validity_end, fuel_included, customs_included, insurance_included,
			fuel_info, customs_info, insurance_info, terms, surcharge_info, is_promotional, is_dynamic,
			data_source, updated_from_offer_id, sample_size, notes, internal_note, meta, meta2, meta3,
			active, deleted
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
			$21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38,
			$39, $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50, $51, $52, $53, 1, 0
		) RETURNING *`

	err := pgxscan.Get(context.Background(), db.DB, &quote, query,
		req.TransportType, req.SubType, req.UserID, req.CompanyID, req.ExecCompanyID, req.VehicleTypeID,
		req.PackagingTypeID, req.CostPerKm, req.Currency, req.FromCountryID, req.FromCityID,
		req.ToCountryID, req.ToCityID, req.Distance, req.FromCountry, req.FromRegion, req.ToCountry, req.ToRegion,
		req.FromAddress, req.ToAddress, req.Tax, req.TaxPrice, req.Trade, req.Discount, req.PaymentMethod, req.PaymentTerm,
		req.DistanceKm, req.AveragePrice, req.MinPrice, req.MaxPrice, req.PriceUnit, req.MinVolume, req.MaxVolume,
		req.ValidityStart, req.ValidityEnd, req.FuelIncluded, req.CustomsIncluded, req.InsuranceIncluded,
		req.FuelInfo, req.CustomsInfo, req.InsuranceInfo, req.Terms, req.SurchargeInfo, req.IsPromotional, req.IsDynamic,
		req.DataSource, req.UpdatedFromOfferID, req.SampleSize, req.Notes, req.InternalNote, req.Meta, req.Meta2, req.Meta3,
	)

	return quote, err
}

func UpdatePriceQuoteRecord(id int, req dto.UpdatePriceQuoteRequest) (dto.PriceQuote, error) {
	var quote dto.PriceQuote

	// Check if exists
	existsQuery := `SELECT id FROM tbl_price_quote WHERE id = $1 AND deleted = 0`
	var existingID int
	err := pgxscan.Get(context.Background(), db.DB, &existingID, existsQuery, id)
	if err != nil {
		return quote, fmt.Errorf("price quote not found")
	}

	setParts := []string{"updated_at = NOW()"}
	args := []interface{}{}
	argIndex := 1

	// Add all fields that can be updated
	if req.TransportType != nil {
		setParts = append(setParts, fmt.Sprintf("transport_type = $%d", argIndex))
		args = append(args, *req.TransportType)
		argIndex++
	}
	if req.SubType != nil {
		setParts = append(setParts, fmt.Sprintf("sub_type = $%d", argIndex))
		args = append(args, *req.SubType)
		argIndex++
	}
	if req.UserID != nil {
		setParts = append(setParts, fmt.Sprintf("user_id = $%d", argIndex))
		args = append(args, *req.UserID)
		argIndex++
	}
	if req.CompanyID != nil {
		setParts = append(setParts, fmt.Sprintf("company_id = $%d", argIndex))
		args = append(args, *req.CompanyID)
		argIndex++
	}
	if req.ExecCompanyID != nil {
		setParts = append(setParts, fmt.Sprintf("exec_company_id = $%d", argIndex))
		args = append(args, *req.ExecCompanyID)
		argIndex++
	}
	if req.VehicleTypeID != nil {
		setParts = append(setParts, fmt.Sprintf("vehicle_type_id = $%d", argIndex))
		args = append(args, *req.VehicleTypeID)
		argIndex++
	}
	if req.PackagingTypeID != nil {
		setParts = append(setParts, fmt.Sprintf("packaging_type_id = $%d", argIndex))
		args = append(args, *req.PackagingTypeID)
		argIndex++
	}
	if req.CostPerKm != nil {
		setParts = append(setParts, fmt.Sprintf("cost_per_km = $%d", argIndex))
		args = append(args, *req.CostPerKm)
		argIndex++
	}
	if req.Currency != nil {
		setParts = append(setParts, fmt.Sprintf("currency = $%d", argIndex))
		args = append(args, *req.Currency)
		argIndex++
	}
	if req.FromCountryID != nil {
		setParts = append(setParts, fmt.Sprintf("from_country_id = $%d", argIndex))
		args = append(args, *req.FromCountryID)
		argIndex++
	}
	if req.FromCityID != nil {
		setParts = append(setParts, fmt.Sprintf("from_city_id = $%d", argIndex))
		args = append(args, *req.FromCityID)
		argIndex++
	}
	if req.ToCountryID != nil {
		setParts = append(setParts, fmt.Sprintf("to_country_id = $%d", argIndex))
		args = append(args, *req.ToCountryID)
		argIndex++
	}
	if req.ToCityID != nil {
		setParts = append(setParts, fmt.Sprintf("to_city_id = $%d", argIndex))
		args = append(args, *req.ToCityID)
		argIndex++
	}
	if req.Distance != nil {
		setParts = append(setParts, fmt.Sprintf("distance = $%d", argIndex))
		args = append(args, *req.Distance)
		argIndex++
	}
	if req.FromCountry != nil {
		setParts = append(setParts, fmt.Sprintf("from_country = $%d", argIndex))
		args = append(args, *req.FromCountry)
		argIndex++
	}
	if req.FromRegion != nil {
		setParts = append(setParts, fmt.Sprintf("from_region = $%d", argIndex))
		args = append(args, *req.FromRegion)
		argIndex++
	}
	if req.ToCountry != nil {
		setParts = append(setParts, fmt.Sprintf("to_country = $%d", argIndex))
		args = append(args, *req.ToCountry)
		argIndex++
	}
	if req.ToRegion != nil {
		setParts = append(setParts, fmt.Sprintf("to_region = $%d", argIndex))
		args = append(args, *req.ToRegion)
		argIndex++
	}
	if req.FromAddress != nil {
		setParts = append(setParts, fmt.Sprintf("from_address = $%d", argIndex))
		args = append(args, *req.FromAddress)
		argIndex++
	}
	if req.ToAddress != nil {
		setParts = append(setParts, fmt.Sprintf("to_address = $%d", argIndex))
		args = append(args, *req.ToAddress)
		argIndex++
	}
	if req.Tax != nil {
		setParts = append(setParts, fmt.Sprintf("tax = $%d", argIndex))
		args = append(args, *req.Tax)
		argIndex++
	}
	if req.TaxPrice != nil {
		setParts = append(setParts, fmt.Sprintf("tax_price = $%d", argIndex))
		args = append(args, *req.TaxPrice)
		argIndex++
	}
	if req.Trade != nil {
		setParts = append(setParts, fmt.Sprintf("trade = $%d", argIndex))
		args = append(args, *req.Trade)
		argIndex++
	}
	if req.Discount != nil {
		setParts = append(setParts, fmt.Sprintf("discount = $%d", argIndex))
		args = append(args, *req.Discount)
		argIndex++
	}
	if req.PaymentMethod != nil {
		setParts = append(setParts, fmt.Sprintf("payment_method = $%d", argIndex))
		args = append(args, *req.PaymentMethod)
		argIndex++
	}
	if req.PaymentTerm != nil {
		setParts = append(setParts, fmt.Sprintf("payment_term = $%d", argIndex))
		args = append(args, *req.PaymentTerm)
		argIndex++
	}
	if req.DistanceKm != nil {
		setParts = append(setParts, fmt.Sprintf("distance_km = $%d", argIndex))
		args = append(args, *req.DistanceKm)
		argIndex++
	}
	if req.AveragePrice != nil {
		setParts = append(setParts, fmt.Sprintf("average_price = $%d", argIndex))
		args = append(args, *req.AveragePrice)
		argIndex++
	}
	if req.MinPrice != nil {
		setParts = append(setParts, fmt.Sprintf("min_price = $%d", argIndex))
		args = append(args, *req.MinPrice)
		argIndex++
	}
	if req.MaxPrice != nil {
		setParts = append(setParts, fmt.Sprintf("max_price = $%d", argIndex))
		args = append(args, *req.MaxPrice)
		argIndex++
	}
	if req.PriceUnit != nil {
		setParts = append(setParts, fmt.Sprintf("price_unit = $%d", argIndex))
		args = append(args, *req.PriceUnit)
		argIndex++
	}
	if req.MinVolume != nil {
		setParts = append(setParts, fmt.Sprintf("min_volume = $%d", argIndex))
		args = append(args, *req.MinVolume)
		argIndex++
	}
	if req.MaxVolume != nil {
		setParts = append(setParts, fmt.Sprintf("max_volume = $%d", argIndex))
		args = append(args, *req.MaxVolume)
		argIndex++
	}
	if req.ValidityStart != nil {
		setParts = append(setParts, fmt.Sprintf("validity_start = $%d", argIndex))
		args = append(args, *req.ValidityStart)
		argIndex++
	}
	if req.ValidityEnd != nil {
		setParts = append(setParts, fmt.Sprintf("validity_end = $%d", argIndex))
		args = append(args, *req.ValidityEnd)
		argIndex++
	}
	if req.FuelIncluded != nil {
		setParts = append(setParts, fmt.Sprintf("fuel_included = $%d", argIndex))
		args = append(args, *req.FuelIncluded)
		argIndex++
	}
	if req.CustomsIncluded != nil {
		setParts = append(setParts, fmt.Sprintf("customs_included = $%d", argIndex))
		args = append(args, *req.CustomsIncluded)
		argIndex++
	}
	if req.InsuranceIncluded != nil {
		setParts = append(setParts, fmt.Sprintf("insurance_included = $%d", argIndex))
		args = append(args, *req.InsuranceIncluded)
		argIndex++
	}
	if req.FuelInfo != nil {
		setParts = append(setParts, fmt.Sprintf("fuel_info = $%d", argIndex))
		args = append(args, *req.FuelInfo)
		argIndex++
	}
	if req.CustomsInfo != nil {
		setParts = append(setParts, fmt.Sprintf("customs_info = $%d", argIndex))
		args = append(args, *req.CustomsInfo)
		argIndex++
	}
	if req.InsuranceInfo != nil {
		setParts = append(setParts, fmt.Sprintf("insurance_info = $%d", argIndex))
		args = append(args, *req.InsuranceInfo)
		argIndex++
	}
	if req.Terms != nil {
		setParts = append(setParts, fmt.Sprintf("terms = $%d", argIndex))
		args = append(args, *req.Terms)
		argIndex++
	}
	if req.SurchargeInfo != nil {
		setParts = append(setParts, fmt.Sprintf("surcharge_info = $%d", argIndex))
		args = append(args, *req.SurchargeInfo)
		argIndex++
	}
	if req.IsPromotional != nil {
		setParts = append(setParts, fmt.Sprintf("is_promotional = $%d", argIndex))
		args = append(args, *req.IsPromotional)
		argIndex++
	}
	if req.IsDynamic != nil {
		setParts = append(setParts, fmt.Sprintf("is_dynamic = $%d", argIndex))
		args = append(args, *req.IsDynamic)
		argIndex++
	}
	if req.DataSource != nil {
		setParts = append(setParts, fmt.Sprintf("data_source = $%d", argIndex))
		args = append(args, *req.DataSource)
		argIndex++
	}
	if req.UpdatedFromOfferID != nil {
		setParts = append(setParts, fmt.Sprintf("updated_from_offer_id = $%d", argIndex))
		args = append(args, *req.UpdatedFromOfferID)
		argIndex++
	}
	if req.SampleSize != nil {
		setParts = append(setParts, fmt.Sprintf("sample_size = $%d", argIndex))
		args = append(args, *req.SampleSize)
		argIndex++
	}
	if req.Notes != nil {
		setParts = append(setParts, fmt.Sprintf("notes = $%d", argIndex))
		args = append(args, *req.Notes)
		argIndex++
	}
	if req.InternalNote != nil {
		setParts = append(setParts, fmt.Sprintf("internal_note = $%d", argIndex))
		args = append(args, *req.InternalNote)
		argIndex++
	}
	if req.Meta != nil {
		setParts = append(setParts, fmt.Sprintf("meta = $%d", argIndex))
		args = append(args, *req.Meta)
		argIndex++
	}
	if req.Meta2 != nil {
		setParts = append(setParts, fmt.Sprintf("meta2 = $%d", argIndex))
		args = append(args, *req.Meta2)
		argIndex++
	}
	if req.Meta3 != nil {
		setParts = append(setParts, fmt.Sprintf("meta3 = $%d", argIndex))
		args = append(args, *req.Meta3)
		argIndex++
	}
	if req.Active != nil {
		setParts = append(setParts, fmt.Sprintf("active = $%d", argIndex))
		args = append(args, *req.Active)
		argIndex++
	}

	args = append(args, id)

	setClause := strings.Join(setParts, ", ")
	query := fmt.Sprintf(`
		UPDATE tbl_price_quote 
		SET %s 
		WHERE id = $%d AND deleted = 0 
		RETURNING *`, setClause, argIndex)

	err = pgxscan.Get(context.Background(), db.DB, &quote, query, args...)
	return quote, err
}

func DeletePriceQuoteRecord(id int) error {
	query := `UPDATE tbl_price_quote SET deleted = 1, updated_at = NOW() WHERE id = $1 AND deleted = 0`
	_, err := db.DB.Exec(context.Background(), query, id)
	return err
}

func GetPriceQuoteWithOfferAnalysis(ctx *gin.Context) {
	var filters dto.PriceQuoteAnalysisFilters
	if err := ctx.ShouldBindQuery(&filters); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid query parameters", err.Error()))
		return
	}

	result, err := AnalyzePriceWithOffers(filters)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to analyze price quotes with offers", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Price analysis completed successfully", result))
}

func AnalyzePriceWithOffers(filters dto.PriceQuoteAnalysisFilters) (dto.PriceQuoteAnalysisResponse, error) {
	var response dto.PriceQuoteAnalysisResponse

	// Step 1: Find matching offers
	offers, offerStats, matchingCriteria, err := findMatchingOffers(filters)
	if err != nil {
		return response, err
	}

	// Step 2: Find matching price quotes as fallback
	priceQuotes, priceQuoteStats, err := findMatchingPriceQuotes(filters)
	if err != nil {
		return response, err
	}

	// Step 3: Build the response
	response.AnalysisInfo.MatchingCriteria = matchingCriteria

	if len(offers) > 0 {
		// Primary source: offers
		response.AnalysisInfo.FoundFromOffers = true
		response.AnalysisInfo.OfferCount = len(offers)
		response.AnalysisInfo.OfferMinPrice = offerStats.MinPrice
		response.AnalysisInfo.OfferMaxPrice = offerStats.MaxPrice
		response.AnalysisInfo.OfferAvgCostPerKm = offerStats.AvgCostPerKm

		// Build PriceQuote structure from offer analysis
		response.PriceQuote = buildPriceQuoteFromOffers(offers, offerStats, filters)
		response.Notes = fmt.Sprintf("Found average among %d offers, minimum price %.2f and maximum %.2f",
			len(offers), offerStats.MinPrice, offerStats.MaxPrice)
	} else if len(priceQuotes) > 0 {
		// Fallback source: price quotes
		response.AnalysisInfo.FoundFromOffers = false
		response.AnalysisInfo.PriceQuoteCount = len(priceQuotes)
		response.AnalysisInfo.PriceQuoteMinPrice = priceQuoteStats.MinPrice
		response.AnalysisInfo.PriceQuoteMaxPrice = priceQuoteStats.MaxPrice
		response.AnalysisInfo.PriceQuoteAvgPrice = priceQuoteStats.AvgPrice

		// Use the best matching price quote
		response.PriceQuote = priceQuotes[0]
		response.Notes = fmt.Sprintf("Found average among %d price quotes, minimum price %.2f and maximum %.2f",
			len(priceQuotes), priceQuoteStats.MinPrice, priceQuoteStats.MaxPrice)
	} else {
		// No matches found
		response.PriceQuote = buildEmptyPriceQuote(filters)
		response.Notes = "No matching offers or price quotes found for the specified criteria"
	}

	return response, nil
}

type OfferStats struct {
	MinPrice     float64
	MaxPrice     float64
	AvgPrice     float64
	AvgCostPerKm float64
}

type PriceQuoteStats struct {
	MinPrice float64
	MaxPrice float64
	AvgPrice float64
}

func findMatchingOffers(filters dto.PriceQuoteAnalysisFilters) ([]map[string]interface{}, OfferStats, []string, error) {
	var offers []map[string]interface{}
	var stats OfferStats
	var matchingCriteria []string

	whereParts := []string{"deleted = 0", "active = 1", "offer_state != 'disabled'"}
	args := []interface{}{}
	argIndex := 1

	if filters.TransportType != "" {
		// Note: We'll need to join with vehicle types or use transport_type from cargo/vehicle data
	}

	if filters.FromCountryID > 0 {
		whereParts = append(whereParts, fmt.Sprintf("from_country_id = $%d", argIndex))
		args = append(args, filters.FromCountryID)
		argIndex++
		matchingCriteria = append(matchingCriteria, "from_country_id")
	} else if filters.FromCountry != "" {
		whereParts = append(whereParts, fmt.Sprintf("from_country ILIKE $%d", argIndex))
		args = append(args, "%"+filters.FromCountry+"%")
		argIndex++
		matchingCriteria = append(matchingCriteria, "from_country")
	}

	if filters.ToCountryID > 0 {
		whereParts = append(whereParts, fmt.Sprintf("to_country_id = $%d", argIndex))
		args = append(args, filters.ToCountryID)
		argIndex++
		matchingCriteria = append(matchingCriteria, "to_country_id")
	} else if filters.ToCountry != "" {
		whereParts = append(whereParts, fmt.Sprintf("to_country ILIKE $%d", argIndex))
		args = append(args, "%"+filters.ToCountry+"%")
		argIndex++
		matchingCriteria = append(matchingCriteria, "to_country")
	}

	if filters.FromCityID > 0 {
		whereParts = append(whereParts, fmt.Sprintf("from_city_id = $%d", argIndex))
		args = append(args, filters.FromCityID)
		argIndex++
		matchingCriteria = append(matchingCriteria, "from_city_id")
	}

	if filters.ToCityID > 0 {
		whereParts = append(whereParts, fmt.Sprintf("to_city_id = $%d", argIndex))
		args = append(args, filters.ToCityID)
		argIndex++
		matchingCriteria = append(matchingCriteria, "to_city_id")
	}

	if filters.VehicleTypeID > 0 {
		whereParts = append(whereParts, fmt.Sprintf("vehicle_type_id = $%d", argIndex))
		args = append(args, filters.VehicleTypeID)
		argIndex++
		matchingCriteria = append(matchingCriteria, "vehicle_type_id")
	}

	if filters.PackagingTypeID > 0 {
		whereParts = append(whereParts, fmt.Sprintf("packaging_type_id = $%d", argIndex))
		args = append(args, filters.PackagingTypeID)
		argIndex++
		matchingCriteria = append(matchingCriteria, "packaging_type_id")
	}

	if filters.Currency != "" {
		whereParts = append(whereParts, fmt.Sprintf("currency = $%d", argIndex))
		args = append(args, filters.Currency)
		argIndex++
		matchingCriteria = append(matchingCriteria, "currency")
	}

	if filters.PaymentMethod != "" {
		whereParts = append(whereParts, fmt.Sprintf("payment_method = $%d", argIndex))
		args = append(args, filters.PaymentMethod)
		argIndex++
		matchingCriteria = append(matchingCriteria, "payment_method")
	}

	if filters.Distance > 0 {
		tolerance := 50 // km tolerance
		if filters.MatchStrict {
			tolerance = 10
		}
		whereParts = append(whereParts, fmt.Sprintf("ABS(distance - $%d) <= $%d", argIndex, argIndex+1))
		args = append(args, filters.Distance, tolerance)
		argIndex += 2
		matchingCriteria = append(matchingCriteria, "distance")
	}

	if !filters.MatchStrict {
		if filters.FromRegion != "" {
			whereParts = append(whereParts, fmt.Sprintf("from_region ILIKE $%d", argIndex))
			args = append(args, "%"+filters.FromRegion+"%")
			argIndex++
			matchingCriteria = append(matchingCriteria, "from_region")
		}

		if filters.ToRegion != "" {
			whereParts = append(whereParts, fmt.Sprintf("to_region ILIKE $%d", argIndex))
			args = append(args, "%"+filters.ToRegion+"%")
			argIndex++
			matchingCriteria = append(matchingCriteria, "to_region")
		}
	}

	whereClause := strings.Join(whereParts, " AND ")

	query := fmt.Sprintf(`
		SELECT 
			id, user_id, company_id, exec_company_id, vehicle_type_id, packaging_type_id,
			cost_per_km, currency, from_country_id, from_city_id, to_country_id, to_city_id,
			distance, from_country, from_region, to_country, to_region, from_address, to_address,
			tax, tax_price, trade, discount, payment_method, payment_term, offer_price, total_price,
			validity_start, validity_end, delivery_start, delivery_end, note, meta, meta2, meta3,
			created_at, updated_at
		FROM tbl_offer 
		WHERE %s 
		ORDER BY created_at DESC 
		LIMIT 100`, whereClause)

	rows, err := db.DB.Query(context.Background(), query, args...)
	if err != nil {
		return offers, stats, matchingCriteria, err
	}
	defer rows.Close()

	var totalPrice, totalCostPerKm float64
	var minPrice, maxPrice float64
	count := 0

	for rows.Next() {
		offer := make(map[string]interface{})
		var id, userID, companyID, execCompanyID, vehicleTypeID, packagingTypeID int
		var costPerKm, taxPrice, offerPrice, totalPriceVal float64
		var currency, fromCountry, fromRegion, toCountry, toRegion, fromAddress, toAddress string
		var paymentMethod, paymentTerm, note, meta, meta2, meta3 string
		var fromCountryID, fromCityID, toCountryID, toCityID, distance, tax, trade, discount int
		var validityStart, validityEnd, deliveryStart, deliveryEnd, createdAt, updatedAt time.Time

		err := rows.Scan(
			&id, &userID, &companyID, &execCompanyID, &vehicleTypeID, &packagingTypeID,
			&costPerKm, &currency, &fromCountryID, &fromCityID, &toCountryID, &toCityID,
			&distance, &fromCountry, &fromRegion, &toCountry, &toRegion, &fromAddress, &toAddress,
			&tax, &taxPrice, &trade, &discount, &paymentMethod, &paymentTerm, &offerPrice, &totalPriceVal,
			&validityStart, &validityEnd, &deliveryStart, &deliveryEnd, &note, &meta, &meta2, &meta3,
			&createdAt, &updatedAt,
		)
		if err != nil {
			continue
		}

		priceToUse := offerPrice
		if priceToUse == 0 {
			priceToUse = totalPriceVal
		}

		if count == 0 {
			minPrice = priceToUse
			maxPrice = priceToUse
		} else {
			if priceToUse < minPrice {
				minPrice = priceToUse
			}
			if priceToUse > maxPrice {
				maxPrice = priceToUse
			}
		}

		totalPrice += priceToUse
		totalCostPerKm += costPerKm
		count++

		offer["id"] = id
		offer["user_id"] = userID
		offer["company_id"] = companyID
		offer["vehicle_type_id"] = vehicleTypeID
		offer["cost_per_km"] = costPerKm
		offer["currency"] = currency
		offer["from_country"] = fromCountry
		offer["to_country"] = toCountry
		offer["distance"] = distance
		offer["offer_price"] = offerPrice
		offer["total_price"] = totalPriceVal
		offer["payment_method"] = paymentMethod
		offer["created_at"] = createdAt

		offers = append(offers, offer)
	}

	if count > 0 {
		stats.MinPrice = minPrice
		stats.MaxPrice = maxPrice
		stats.AvgPrice = totalPrice / float64(count)
		stats.AvgCostPerKm = totalCostPerKm / float64(count)
	}

	return offers, stats, matchingCriteria, nil
}

func findMatchingPriceQuotes(filters dto.PriceQuoteAnalysisFilters) ([]dto.PriceQuote, PriceQuoteStats, error) {
	var quotes []dto.PriceQuote
	var stats PriceQuoteStats

	priceQuoteFilters := dto.PriceQuoteFilters{
		TransportType: filters.TransportType,
		SubType:       filters.SubType,
		FromCountry:   filters.FromCountry,
		ToCountry:     filters.ToCountry,
		FromRegion:    filters.FromRegion,
		ToRegion:      filters.ToRegion,
		Currency:      filters.Currency,
		VehicleTypeID: filters.VehicleTypeID,
		PaymentMethod: filters.PaymentMethod,
		Page:          1,
		PerPage:       50,
		SortBy:        "created_at",
		SortOrder:     "DESC",
	}

	quotes, _, err := GetPriceQuotes(priceQuoteFilters)
	if err != nil {
		return quotes, stats, err
	}

	if len(quotes) > 0 {
		var total float64
		minPrice := quotes[0].AveragePrice
		maxPrice := quotes[0].AveragePrice

		for _, quote := range quotes {
			total += quote.AveragePrice
			if quote.AveragePrice < minPrice {
				minPrice = quote.AveragePrice
			}
			if quote.AveragePrice > maxPrice {
				maxPrice = quote.AveragePrice
			}
		}

		stats.MinPrice = minPrice
		stats.MaxPrice = maxPrice
		stats.AvgPrice = total / float64(len(quotes))
	}

	return quotes, stats, nil
}

func buildPriceQuoteFromOffers(offers []map[string]interface{}, stats OfferStats, filters dto.PriceQuoteAnalysisFilters) dto.PriceQuote {
	if len(offers) == 0 {
		return buildEmptyPriceQuote(filters)
	}

	firstOffer := offers[0]

	quote := dto.PriceQuote{
		TransportType:     filters.TransportType,
		SubType:           filters.SubType,
		Currency:          getString(firstOffer, "currency"),
		FromCountryID:     getInt(firstOffer, "from_country_id"),
		FromCityID:        getInt(firstOffer, "from_city_id"),
		ToCountryID:       getInt(firstOffer, "to_country_id"),
		ToCityID:          getInt(firstOffer, "to_city_id"),
		FromCountry:       getString(firstOffer, "from_country"),
		ToCountry:         getString(firstOffer, "to_country"),
		Distance:          getInt(firstOffer, "distance"),
		VehicleTypeID:     getInt(firstOffer, "vehicle_type_id"),
		PackagingTypeID:   getInt(firstOffer, "packaging_type_id"),
		CostPerKm:         stats.AvgCostPerKm,
		AveragePrice:      stats.AvgPrice,
		MinPrice:          stats.MinPrice,
		MaxPrice:          stats.MaxPrice,
		PriceUnit:         "per_trip", // Default assumption
		PaymentMethod:     getString(firstOffer, "payment_method"),
		ValidityStart:     time.Now(),
		ValidityEnd:       time.Now().AddDate(0, 3, 0), // 3 months validity
		DataSource:        "offer_based",
		SampleSize:        len(offers),
		IsPromotional:     false,
		IsDynamic:         true,
		FuelIncluded:      true, // Default assumption
		CustomsIncluded:   false,
		InsuranceIncluded: false,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		Active:            1,
		Deleted:           0,
	}

	return quote
}

func buildEmptyPriceQuote(filters dto.PriceQuoteAnalysisFilters) dto.PriceQuote {
	return dto.PriceQuote{
		TransportType:   filters.TransportType,
		SubType:         filters.SubType,
		Currency:        filters.Currency,
		FromCountryID:   filters.FromCountryID,
		FromCityID:      filters.FromCityID,
		ToCountryID:     filters.ToCountryID,
		ToCityID:        filters.ToCityID,
		FromCountry:     filters.FromCountry,
		ToCountry:       filters.ToCountry,
		VehicleTypeID:   filters.VehicleTypeID,
		PackagingTypeID: filters.PackagingTypeID,
		Distance:        filters.Distance,
		DistanceKm:      filters.DistanceKm,
		PaymentMethod:   filters.PaymentMethod,
		MinVolume:       filters.MinVolume,
		MaxVolume:       filters.MaxVolume,
		AveragePrice:    0,
		MinPrice:        0,
		MaxPrice:        0,
		PriceUnit:       "unknown",
		ValidityStart:   time.Now(),
		ValidityEnd:     time.Now().AddDate(1, 0, 0),
		DataSource:      "no_match",
		IsPromotional:   false,
		IsDynamic:       false,
		FuelIncluded:    filters.FuelIncluded != nil && *filters.FuelIncluded,
		CustomsIncluded: filters.CustomsIncluded != nil && *filters.CustomsIncluded,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Active:          1,
		Deleted:         0,
	}
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getInt(m map[string]interface{}, key string) int {
	if v, ok := m[key]; ok {
		if i, ok := v.(int); ok {
			return i
		}
	}
	return 0
}
