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

	// TODO: I like this code approach
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
			// TODO: but needs sql fix
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
			&offer.ExecCompanyID, &offer.DriverID, &offer.VehicleID, &offer.VehicleTypeID,
			&offer.CargoID, &offer.PackagingTypeID, &offer.OfferState, &offer.OfferRole,
			&offer.CostPerKm, &offer.Currency, &offer.FromCountryID,
			&offer.FromCityID, &offer.ToCountryID, &offer.ToCityID, &offer.Distance,
			&offer.FromCountry, &offer.FromRegion, &offer.ToCountry,
			&offer.ToRegion, &offer.FromAddress, &offer.ToAddress, &offer.MapURL,
			&offer.SenderContact, &offer.RecipientContact,
			&offer.DeliverContact, &offer.ViewCount, &offer.ValidityStart,
			&offer.ValidityEnd, &offer.DeliveryStart, &offer.DeliveryEnd,
			&offer.Note, &offer.Tax, &offer.TaxPrice, &offer.Trade,
			&offer.Discount, &offer.PaymentMethod, &offer.PaymentTerm, &offer.Meta,
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

func CreateOffer(ctx *gin.Context) {
	var offer dto.Offer
	if err := ctx.ShouldBindJSON(&offer); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	companyID := ctx.MustGet("companyID").(int)
	userID := ctx.MustGet("id").(int)
	role := ctx.MustGet("role").(string)
	if !(role == "admin" || role == "system") {
		offer.CompanyID = companyID
		offer.UserID = userID
		offer.OfferRole = role
	}

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
		offer.VehicleTypeID, offer.PackagingTypeID, offer.Distance, offer.MapURL, offer.PaymentTerm,
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
		offer.Active = nil
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
		offer.OfferState, offer.OfferRole,
		offer.VehicleTypeID, offer.PackagingTypeID, offer.Distance, offer.MapURL, offer.PaymentTerm,
		offer.CompanyID,
	).Scan(&updatedID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error updating offer", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.FormatResponse("Successfully updated offer!", gin.H{"id": updatedID}))
}

func GetOffer(ctx *gin.Context) {
	const query = `
        SELECT 
            o.*,
            json_build_object(
                'id', c.id,
                'company_name', c.company_name,
                'email', c.email,
                'phone', c.phone,
                'address', c.address,
                'country', c.country
            ) as company,
            json_build_object(
                'id', ec.id,
                'company_name', ec.company_name,
                'email', ec.email,
                'phone', ec.phone,
                'address', ec.address,
                'country', ec.country
            ) as exec_company,
            json_build_object(
                'id', d.id,
                'first_name', d.first_name,
                'last_name', d.last_name,
                'phone', d.phone,
                'email', d.email,
                'image_url', d.image_url
            ) as assigned_driver,
            json_build_object(
                'id', v.id,
                'vehicle_type_id', v.vehicle_type_id,
                'numberplate', v.numberplate,
                'mileage', v.mileage,
                'gps', v.gps,
                'available', v.available
            ) as vehicle,
            json_build_object(
                'id', vt.id,
                'title_en', vt.title_en,
                'title_ru', vt.title_ru,
                'title_tk', vt.title_tk
            ) as vehicle_type,
            json_build_object(
                'id', cg.id,
                'name', cg.name,
                'description', cg.description,
                'info', cg.info
            ) as cargo,
            json_build_object(
                'id', pt.id,
                'name_en', pt.name_en,
                'name_ru', pt.name_ru,
                'name_tk', pt.name_tk
            ) as packaging_type
        FROM tbl_offer o
        LEFT JOIN tbl_company c ON o.company_id = c.id
        LEFT JOIN tbl_company ec ON o.exec_company_id = ec.id
        LEFT JOIN tbl_driver d ON o.driver_id = d.id
        LEFT JOIN tbl_vehicle v ON o.vehicle_id = v.id
        LEFT JOIN tbl_vehicle_type vt ON o.vehicle_type_id = vt.id
        LEFT JOIN tbl_cargo cg ON o.cargo_id = cg.id
        LEFT JOIN tbl_packaging_type pt ON o.packaging_type_id = pt.id
        WHERE o.id = $1 AND o.deleted = 0`

	var offer dto.OfferDetailedResponse
	err := pgxscan.Get(context.Background(), db.DB, &offer, query, ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Offer not found", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Offer details", offer))
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

func GetDetailedOfferList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
	offset := (page - 1) * perPage

	validOrderColumns := map[string]bool{
		"id": true, "cost_per_km": true, "distance": true,
		"view_count": true, "created_at": true, "updated_at": true,
	}

	orderBy := ctx.DefaultQuery("order_by", "id")
	if !validOrderColumns[orderBy] {
		orderBy = "id"
	}
	orderDir := strings.ToUpper(ctx.DefaultQuery("order_dir", "DESC"))
	if orderDir != "ASC" && orderDir != "DESC" {
		orderDir = "DESC"
	}

	baseQuery := `
        WITH company_stats AS (
            SELECT 
                c.id as company_id,
                COUNT(DISTINCT d.id) as drivers_count,
                COUNT(DISTINCT CASE 
                    WHEN o2.active = 1 
                    AND o2.deleted = 0 
                    AND o2.validity_end >= CURRENT_DATE 
                    THEN o2.id 
                    END) as offers_count
            FROM tbl_company c
            LEFT JOIN tbl_driver d ON d.company_id = c.id AND d.deleted = 0 AND d.active = 1
            LEFT JOIN tbl_offer o2 ON (o2.company_id = c.id OR o2.exec_company_id = c.id)
            GROUP BY c.id
        )
        SELECT 
            o.*,
            COUNT(*) OVER() as total_count,
            
            -- Company fields with stats
            json_build_object(
				'id', c.id,
				'uuid', c.uuid,
				'user_id', c.user_id,
				'role_id', c.role_id,
				'company_name', c.company_name,
				'first_name', c.first_name,
				'last_name', c.last_name,
				'patronymic_name', c.patronymic_name,
				'phone', c.phone,
				'phone2', c.phone2,
				'phone3', c.phone3,
				'email', c.email,
				'email2', c.email2,
				'email3', c.email3,
				'meta', c.meta,
				'meta2', c.meta2,
				'meta3', c.meta3,
				'address', c.address,
				'country', c.country,
				'country_id', c.country_id,
				'city_id', c.city_id,
				'image_url', c.image_url,
				'entity', c.entity,
				'featured', c.featured,
				'rating', c.rating,
				'partner', c.partner,
                'drivers_count', COALESCE(cs1.drivers_count, 0),
                'offers_count', COALESCE(cs1.offers_count, 0)
            ) as company,
            
            -- Exec Company fields with stats
            json_build_object(
				'id', ec.id,
				'uuid', ec.uuid,
				'user_id', ec.user_id,
				'role_id', ec.role_id,
				'company_name', ec.company_name,
				'first_name', ec.first_name,
				'last_name', ec.last_name,
				'patronymic_name', ec.patronymic_name,
				'phone', ec.phone,
				'phone2', ec.phone2,
				'phone3', ec.phone3,
				'email', ec.email,
				'email2', ec.email2,
				'email3', ec.email3,
				'meta', ec.meta,
				'meta2', ec.meta2,
				'meta3', ec.meta3,
				'address', ec.address,
				'country', ec.country,
				'country_id', ec.country_id,
				'city_id', ec.city_id,
				'image_url', ec.image_url,
				'entity', ec.entity,
				'featured', ec.featured,
				'rating', ec.rating,
				'partner', ec.partner,
                'drivers_count', COALESCE(cs2.drivers_count, 0),
                'offers_count', COALESCE(cs2.offers_count, 0)
            ) as exec_company,
            
            -- Driver fields
            json_build_object(
				'id',d.id,
				'uuid',d.uuid,
				'company_id',d.company_id,
				'first_name',d.first_name,
				'last_name',d.last_name,
				'patronymic_name',d.patronymic_name,
				'phone',d.phone,
				'email',d.email,
				'featured',d.featured,
				'rating',d.rating,
				'partner',d.partner,
				'successful_ops',d.successful_ops,
				'image_url',d.image_url,
				'view_count',d.view_count,
				'meta',d.meta,
				'meta2',d.meta2,
				'meta3',d.meta3,
				'available',d.available
            ) as assigned_driver,
            
            -- Vehicle fields
            json_build_object(
                'id', v.id,
                'vehicle_type_id', v.vehicle_type_id,
                'numberplate', v.numberplate,
                'mileage', v.mileage,
                'gps', v.gps,
                'available', v.available
            ) as vehicle,
            
            -- Vehicle type fields
            json_build_object(
                'id', vt.id,
                'title_en', vt.title_en,
                'title_ru', vt.title_ru,
                'title_tk', vt.title_tk
            ) as vehicle_type,
            
            -- Cargo fields
            json_build_object(
                'id', cg.id,
                'name', cg.name,
                'description', cg.description,
                'info', cg.info,
                'qty', cg.qty,
                'weight', cg.weight
            ) as cargo,
            
            -- Packaging type fields
            json_build_object(
                'id', pt.id,
                'name_en', pt.name_en,
                'name_ru', pt.name_ru,
                'name_tk', pt.name_tk,
                'material', pt.material,
                'dimensions', pt.dimensions
            ) as packaging_type
            
        FROM tbl_offer o
        LEFT JOIN tbl_company c ON o.company_id = c.id
        LEFT JOIN company_stats cs1 ON c.id = cs1.company_id
        LEFT JOIN tbl_company ec ON o.exec_company_id = ec.id
        LEFT JOIN company_stats cs2 ON ec.id = cs2.company_id
        LEFT JOIN tbl_driver d ON o.driver_id = d.id
        LEFT JOIN tbl_vehicle v ON o.vehicle_id = v.id
        LEFT JOIN tbl_vehicle_type vt ON o.vehicle_type_id = vt.id
        LEFT JOIN tbl_cargo cg ON o.cargo_id = cg.id
        LEFT JOIN tbl_packaging_type pt ON o.packaging_type_id = pt.id
        `

	var whereClauses []string
	var args []interface{}
	argCounter := 1

	role := ctx.MustGet("role").(string)
	if !(role == "admin" || role == "system") {
		whereClauses = append(whereClauses, "o.deleted = 0")
		whereClauses = append(whereClauses, "o.active = 1")
	}

	filters := map[string]string{
		"company_id":        ctx.Query("company_id"),
		"exec_company_id":   ctx.Query("exec_company_id"),
		"driver_id":         ctx.Query("driver_id"),
		"vehicle_id":        ctx.Query("vehicle_id"),
		"vehicle_type_id":   ctx.Query("vehicle_type_id"),
		"cargo_id":          ctx.Query("cargo_id"),
		"packaging_type_id": ctx.Query("packaging_type_id"),
		"offer_state":       ctx.Query("offer_state"),
		"offer_role":        ctx.Query("offer_role"),
		"currency":          ctx.Query("currency"),
		"payment_method":    ctx.Query("payment_method"),
	}

	for key, value := range filters {
		if value != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("o.%s = $%d", key, argCounter))
			args = append(args, value)
			argCounter++
		}
	}

	numericRanges := map[string]struct {
		min string
		max string
	}{
		"cost_per_km": {ctx.Query("min_cost"), ctx.Query("max_cost")},
		"distance":    {ctx.Query("min_distance"), ctx.Query("max_distance")},
		"tax":         {ctx.Query("min_tax"), ctx.Query("max_tax")},
		"discount":    {ctx.Query("min_discount"), ctx.Query("max_discount")},
	}

	for field, ranges := range numericRanges {
		if ranges.min != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("o.%s >= $%d", field, argCounter))
			minVal, _ := strconv.ParseFloat(ranges.min, 64)
			args = append(args, minVal)
			argCounter++
		}
		if ranges.max != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("o.%s <= $%d", field, argCounter))
			maxVal, _ := strconv.ParseFloat(ranges.max, 64)
			args = append(args, maxVal)
			argCounter++
		}
	}

	dateRanges := map[string]struct {
		start string
		end   string
	}{
		"validity": {ctx.Query("validity_start"), ctx.Query("validity_end")},
		"delivery": {ctx.Query("delivery_start"), ctx.Query("delivery_end")},
	}

	for field, ranges := range dateRanges {
		if ranges.start != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("o.%s_start >= $%d", field, argCounter))
			args = append(args, ranges.start)
			argCounter++
		}
		if ranges.end != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("o.%s_end <= $%d", field, argCounter))
			args = append(args, ranges.end)
			argCounter++
		}
	}

	locationFilters := map[string]string{
		"from_country_id": ctx.Query("from_country_id"),
		"from_city_id":    ctx.Query("from_city_id"),
		"to_country_id":   ctx.Query("to_country_id"),
		"to_city_id":      ctx.Query("to_city_id"),
	}

	for key, value := range locationFilters {
		if value != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("o.%s = $%d", key, argCounter))
			args = append(args, value)
			argCounter++
		}
	}

	searchTerm := ctx.Query("search")
	if searchTerm != "" {
		searchClause := fmt.Sprintf(`(
                    o.note ILIKE $%d OR 
                    o.sender_contact ILIKE $%d OR
                    o.recipient_contact ILIKE $%d OR
                    o.deliver_contact ILIKE $%d
                )`, argCounter, argCounter, argCounter, argCounter, argCounter, argCounter)
		whereClauses = append(whereClauses, searchClause)
		args = append(args, "%"+searchTerm+"%")
		argCounter++
	}

	searchFromLocation := ctx.Query("from_location")
	if searchFromLocation != "" {
		searchClause := fmt.Sprintf(`(
                    o.from_country ILIKE $%d OR
                    o.from_region ILIKE $%d OR 
                    o.from_address ILIKE $%d OR 
                )`, argCounter, argCounter, argCounter, argCounter, argCounter, argCounter)
		whereClauses = append(whereClauses, searchClause)
		args = append(args, "%"+searchFromLocation+"%")
		argCounter++
	}
	searchToLocation := ctx.Query("to_location")
	if searchToLocation != "" {
		searchClause := fmt.Sprintf(`(
                    o.to_country ILIKE $%d OR
                    o.to_region ILIKE $%d OR 
                    o.to_address ILIKE $%d OR 
                )`, argCounter, argCounter, argCounter, argCounter, argCounter, argCounter)
		whereClauses = append(whereClauses, searchClause)
		args = append(args, "%"+searchToLocation+"%")
		argCounter++
	}

	query := baseQuery
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	query += fmt.Sprintf(" ORDER BY o.%s %s LIMIT $%d OFFSET $%d",
		orderBy, orderDir, argCounter, argCounter+1)
	args = append(args, perPage, offset)

	var offers []dto.OfferDetailedResponse
	err := pgxscan.Select(
		context.Background(),
		db.DB,
		&offers,
		query,
		args...,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			utils.FormatErrorResponse("Couldn't retrieve data", err.Error()))
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

	ctx.JSON(http.StatusOK, utils.FormatResponse("Offer list detailed", response))
}
