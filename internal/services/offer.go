package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	db "texApi/database"
	"texApi/internal/dto"
	"texApi/internal/queries"
	"texApi/pkg/utils"
	"time"
)

func GetMyOfferListUpdate(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
	offset := (page - 1) * perPage

	filters := map[string]interface{}{
		"o.company_id":      ctx.MustGet("companyID"),
		"o.exec_company_id": ctx.Query("exec_company_id"),
		"o.driver_id":       ctx.Query("driver_id"),
		"o.vehicle_id":      ctx.Query("vehicle_id"),
		"o.cargo_id":        ctx.Query("cargo_id"),
		"o.offer_state":     ctx.Query("offer_state"),
		"o.offer_role":      ctx.Query("offer_role"),
		"o.from_country_id": ctx.Query("from_country_id"),
		"o.from_city_id":    ctx.Query("from_city_id"),
		"o.to_country_id":   ctx.Query("to_country_id"),
		"o.to_city_id":      ctx.Query("to_city_id"),
		"o.tax":             ctx.Query("tax"),
		"o.trade":           ctx.Query("trade"),
		"o.discount":        ctx.Query("discount"),
		"o.payment_method":  ctx.Query("payment_method"),
		"o.featured":        ctx.Query("featured"),
		"o.partner":         ctx.Query("partner"),
		"o.active":          ctx.Query("active"),
	}

	validityStart := ctx.Query("validity_start")
	validityEnd := ctx.Query("validity_end")

	orderBy := ctx.DefaultQuery("order_by", "id")
	orderDir := ctx.DefaultQuery("order_dir", "DESC")

	var whereClauses []string
	var args []interface{}
	argCounter := 1

	for key, value := range filters {
		if value != "" && value != nil {
			whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", key, argCounter))
			args = append(args, value)
			argCounter++
		}
	}

	if validityStart != "" && validityEnd != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("validity_start BETWEEN $%d AND $%d", argCounter, argCounter+1))
		startTime, _ := time.Parse(time.RFC3339, validityStart)
		endTime, _ := time.Parse(time.RFC3339, validityEnd)
		args = append(args, startTime, endTime)
		argCounter += 2
	} else if validityStart != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("validity_start >= $%d", argCounter))
		startTime, _ := time.Parse(time.RFC3339, validityStart)
		args = append(args, startTime)
		argCounter++
	} else if validityEnd != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("validity_start <= $%d", argCounter))
		endTime, _ := time.Parse(time.RFC3339, validityEnd)
		args = append(args, endTime)
		argCounter++
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	query := fmt.Sprintf(`
        SELECT 
			o.*, 
			COUNT(*) OVER () as total_count,
			c as company,
			d as assigned_driver,
			v as assigned_vehicle,
			cr as cargo
		FROM tbl_offer o
		LEFT JOIN tbl_company c ON o.company_id = c.id
		LEFT JOIN tbl_driver d ON o.driver_id = d.id
		LEFT JOIN tbl_vehicle v ON o.vehicle_id = v.id
		LEFT JOIN tbl_cargo cr ON o.cargo_id = cr.id
		%s
		ORDER BY o.%s %s
		LIMIT $%d OFFSET $%d
    `, whereClause, orderBy, orderDir, argCounter, argCounter+1)

	args = append(args, perPage, offset)

	rows, err := db.DB.Query(context.Background(), query, args...)
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
			&offer.ID, &offer.UUID, &offer.UserID, &offer.CompanyID,
			&offer.ExecCompanyID, &offer.DriverID, &offer.VehicleID,
			&offer.CargoID, &offer.OfferState, &offer.OfferRole,
			&offer.CostPerKm, &offer.Currency, &offer.FromCountryID,
			&offer.FromCityID, &offer.ToCountryID, &offer.ToCityID,
			&offer.FromCountry, &offer.FromRegion, &offer.ToCountry,
			&offer.ToRegion, &offer.FromAddress, &offer.ToAddress,
			&offer.SenderContact, &offer.RecipientContact,
			&offer.DeliverContact, &offer.ViewCount, &offer.ValidityStart,
			&offer.ValidityEnd, &offer.DeliveryStart, &offer.DeliveryEnd,
			&offer.Note, &offer.Tax, &offer.TaxPrice, &offer.Trade,
			&offer.Discount, &offer.PaymentMethod, &offer.Meta,
			&offer.Meta2, &offer.Meta3, &offer.Featured, &offer.Partner,
			&offer.CreatedAt, &offer.UpdatedAt, &offer.Active,
			&offer.Deleted, &totalCount,
			&companyJSON, &driverJSON, &vehicleJSON, &cargoJSON,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Scan error", err.Error()))
			return
		}

		json.Unmarshal(companyJSON, &offer.Company)
		json.Unmarshal(driverJSON, &offer.AssignedDriver)
		json.Unmarshal(vehicleJSON, &offer.AssignedVehicle)
		// Uncomment if cargo unmarshaling is needed
		json.Unmarshal(cargoJSON, &offer.Cargo)

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

func GetMyOfferList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
	offset := (page - 1) * perPage

	companyID := ctx.MustGet("companyID").(int)
	offerState := ctx.GetHeader("OfferState")

	rows, err := db.DB.Query(
		context.Background(),
		queries.GetMyOfferList,
		companyID,
		perPage,
		offset,
		offerState,
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
			&offer.ID, &offer.UUID, &offer.UserID, &offer.CompanyID, &offer.ExecCompanyID, &offer.DriverID, &offer.VehicleID, &offer.CargoID,
			&offer.OfferState, &offer.OfferRole, &offer.CostPerKm, &offer.Currency,
			&offer.FromCountryID, &offer.FromCityID, &offer.ToCountryID, &offer.ToCityID,
			&offer.FromCountry, &offer.FromRegion, &offer.ToCountry, &offer.ToRegion,
			&offer.FromAddress, &offer.ToAddress, &offer.SenderContact, &offer.RecipientContact, &offer.DeliverContact,
			&offer.ViewCount, &offer.ValidityStart, &offer.ValidityEnd, &offer.DeliveryStart, &offer.DeliveryEnd,
			&offer.Note, &offer.Tax, &offer.TaxPrice, &offer.Trade, &offer.Discount, &offer.PaymentMethod, &offer.Meta, &offer.Meta2, &offer.Meta3,
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

func GetOfferListUpdate(ctx *gin.Context) {

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
	offset := (page - 1) * perPage

	filters := map[string]interface{}{
		"company_id":      ctx.Query("company_id"),
		"exec_company_id": ctx.Query("exec_company_id"),
		"driver_id":       ctx.Query("driver_id"),
		"vehicle_id":      ctx.Query("vehicle_id"),
		"cargo_id":        ctx.Query("cargo_id"),
		"offer_state":     ctx.Query("offer_state"),
		"offer_role":      ctx.Query("offer_role"),
		"from_country_id": ctx.Query("from_country_id"),
		"from_city_id":    ctx.Query("from_city_id"),
		"to_country_id":   ctx.Query("to_country_id"),
		"to_city_id":      ctx.Query("to_city_id"),
		"tax":             ctx.Query("tax"),
		"trade":           ctx.Query("trade"),
		"discount":        ctx.Query("discount"),
		"payment_method":  ctx.Query("payment_method"),
		"featured":        ctx.Query("featured"),
		"partner":         ctx.Query("partner"),
		"active":          ctx.Query("active"),
	}

	validityStart := ctx.Query("validity_start")
	validityEnd := ctx.Query("validity_end")
	deliveryStart := ctx.Query("delivery_start")
	deliveryEnd := ctx.Query("delivery_end")

	orderBy := ctx.DefaultQuery("order_by", "o.id")
	orderDir := ctx.DefaultQuery("order_dir", "DESC")

	role := ctx.MustGet("role").(string)

	stmt := `
        SELECT 
            o.*,
            COUNT(*) OVER() as total_count
        FROM tbl_offer o
        WHERE 
            o.validity_end > CURRENT_TIMESTAMP
            AND o.delivery_end > CURRENT_TIMESTAMP
    `

	var whereClauses []string
	var args []interface{}
	argCounter := 1

	if !(role == "admin" || role == "system") {
		whereClauses = append(whereClauses, "o.offer_state = 'enabled' AND o.deleted = 0")
	}

	for key, value := range filters {
		if value != "" && value != nil {
			whereClauses = append(whereClauses, fmt.Sprintf("o.%s = $%d", key, argCounter))
			args = append(args, value)
			argCounter++
		}
	}

	if validityStart != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("o.validity_start >= $%d", argCounter))
		startTime, _ := time.Parse(time.RFC3339, validityStart)
		args = append(args, startTime)
		argCounter++
	}
	if validityEnd != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("o.validity_end <= $%d", argCounter))
		endTime, _ := time.Parse(time.RFC3339, validityEnd)
		args = append(args, endTime)
		argCounter++
	}

	if deliveryStart != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("o.delivery_start >= $%d", argCounter))
		startTime, _ := time.Parse(time.RFC3339, deliveryStart)
		args = append(args, startTime)
		argCounter++
	}
	if deliveryEnd != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("o.delivery_end <= $%d", argCounter))
		endTime, _ := time.Parse(time.RFC3339, deliveryEnd)
		args = append(args, endTime)
		argCounter++
	}

	if len(whereClauses) > 0 {
		stmt += " AND " + strings.Join(whereClauses, " AND ")
	}

	stmt += fmt.Sprintf(" ORDER BY %s %s", orderBy, orderDir)
	stmt += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCounter, argCounter+1)
	args = append(args, perPage, offset)

	var offers []dto.Offer
	err := pgxscan.Select(
		context.Background(),
		db.DB,
		&offers,
		stmt,
		args...,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Couldn't retrieve data", err.Error()))
		return
	}

	var totalCount int
	if len(offers) > 0 {
		totalCount = offers[0].TotalCount
	}

	response := utils.PaginatedResponse{
		Total:   totalCount,
		Page:    page,
		PerPage: perPage,
		Data:    offers,
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Offer list", response))
}

func GetOfferList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
	offset := (page - 1) * perPage

	// TODO MAKE DELETED SEPARATED
	companyID, _ := strconv.Atoi(ctx.GetHeader("CompanyID"))
	stmt := queries.GetOfferList
	role := ctx.MustGet("role").(string)
	if !(role == "admin" || role == "system") {
		stmt += ` AND o.offer_state='enabled'`
		stmt += ` AND (o.company_id = $3 OR $3 = 0) AND o.deleted = 0`
	} else {
		stmt += ` AND (o.company_id = $3 OR $3 = 0)`
	}
	stmt += ` ORDER BY o.id DESC LIMIT $1 OFFSET $2;`

	var offers []dto.Offer

	err := pgxscan.Select(
		context.Background(),
		db.DB,
		&offers,
		stmt,
		perPage,
		offset,
		companyID,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Couldn't retrieve data", err.Error()))
		return
	}

	var totalCount int
	if len(offers) > 0 {
		totalCount = offers[0].TotalCount
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

	var offer dto.OfferDetails
	var companyJSON, driverJSON, vehicleJSON, cargoJSON []byte

	err := db.DB.QueryRow(
		context.Background(),
		queries.GetOfferByID,
		id,
	).Scan(
		&offer.ID, &offer.UUID, &offer.UserID, &offer.CompanyID, &offer.ExecCompanyID, &offer.DriverID, &offer.VehicleID, &offer.CargoID,
		&offer.OfferState, &offer.OfferRole, &offer.CostPerKm, &offer.Currency,
		&offer.FromCountryID, &offer.FromCityID, &offer.ToCountryID, &offer.ToCityID,
		&offer.FromCountry, &offer.FromRegion, &offer.ToCountry, &offer.ToRegion,
		&offer.FromAddress, &offer.ToAddress, &offer.SenderContact, &offer.RecipientContact, &offer.DeliverContact,
		&offer.ViewCount, &offer.ValidityStart, &offer.ValidityEnd, &offer.DeliveryStart, &offer.DeliveryEnd,
		&offer.Note, &offer.Tax, &offer.TaxPrice, &offer.Trade, &offer.Discount, &offer.PaymentMethod, &offer.Meta, &offer.Meta2, &offer.Meta3,
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
	json.Unmarshal(cargoJSON, &offer.Cargo)

	ctx.JSON(http.StatusOK, utils.FormatResponse("Offer details", offer))
}

func CreateOffer(ctx *gin.Context) {
	var offer dto.Offer
	if err := ctx.ShouldBindJSON(&offer); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	// TODO MAKE ADMIN STUFFF HERE
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
		offer.FromCountryID, offer.FromCityID, offer.ToCountryID, offer.ToCityID,
		offer.FromCountry, offer.FromRegion, offer.ToCountry, offer.ToRegion, offer.FromAddress, offer.ToAddress,
		offer.SenderContact, offer.RecipientContact, offer.DeliverContact, offer.ValidityStart, offer.ValidityEnd,
		offer.DeliveryStart, offer.DeliveryEnd, offer.Note, offer.Tax, offer.TaxPrice, offer.Trade, offer.Discount,
		offer.PaymentMethod, offer.Meta, offer.Meta2, offer.Meta3, offer.OfferRole, offer.ExecCompanyID,
	).Scan(&id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating offer", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.FormatResponse("Successfully created offer!", gin.H{"id": id}))
}

func UpdateOffer(ctx *gin.Context) {
	offerID, _ := strconv.Atoi(ctx.Param("id"))
	var offer dto.OfferUpdate
	stmt := queries.UpdateOffer

	if err := ctx.ShouldBindJSON(&offer); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	role := ctx.MustGet("role").(string)
	if !(role == "admin" || role == "system") {
		companyID := ctx.MustGet("companyID").(int)
		offer.CompanyID = &companyID
		offer.OfferState = nil
		offer.OfferRole = nil
		stmt += ` AND active = 1 AND deleted = 0`
	}
	stmt += ` RETURNING id;`

	var updatedID int
	err := db.DB.QueryRow(
		context.Background(),
		stmt,
		offerID,
		offer.DriverID, offer.VehicleID, offer.CargoID, offer.CostPerKm, offer.Currency,
		offer.FromCountryID, offer.FromCityID, offer.ToCountryID, offer.ToCityID,
		offer.FromCountry, offer.FromRegion, offer.ToCountry, offer.ToRegion,
		offer.FromAddress, offer.ToAddress,
		offer.SenderContact, offer.RecipientContact, offer.DeliverContact,
		offer.ValidityStart, offer.ValidityEnd, offer.DeliveryStart, offer.DeliveryEnd,
		offer.Note, offer.Tax, offer.TaxPrice, offer.Trade, offer.Discount,
		offer.PaymentMethod, offer.Meta, offer.Meta2, offer.Meta3,
		offer.Active, offer.Deleted, offer.ExecCompanyID,
		offer.OfferState, offer.OfferRole, offer.CompanyID,
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

	result, err := db.DB.Exec(
		context.Background(),
		queries.DeleteOffer,
		id, companyID,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error deleting offer", err.Error()))
		return
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Offer not found", "No offer found with the given ID for this company"))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully deleted offer!", gin.H{"id": id}))
}
