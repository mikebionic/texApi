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
