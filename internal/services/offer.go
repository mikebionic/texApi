package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
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

	// Valid order columns
	validOrderColumns := map[string]bool{
		"id": true, "cost_per_km": true, "distance": true,
		"offer_price": true, "total_price": true, "view_count": true,
		"validity_start": true, "validity_end": true, "delivery_start": true,
		"delivery_end": true, "created_at": true, "updated_at": true,
		"tax": true, "discount": true,
	}

	orderBy := ctx.DefaultQuery("order_by", "id")
	if !validOrderColumns[orderBy] {
		orderBy = "id"
	}
	orderDir := strings.ToUpper(ctx.DefaultQuery("order_dir", "DESC"))
	if orderDir != "ASC" && orderDir != "DESC" {
		orderDir = "DESC"
	}

	filters := map[string]interface{}{
		"o.company_id":      ctx.MustGet("companyID"),
		"o.exec_company_id": ctx.Query("exec_company_id"),
		"o.driver_id":       ctx.Query("driver_id"),
		"o.vehicle_id":      ctx.Query("vehicle_id"),
		"o.trailer_id":      ctx.Query("trailer_id"),
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
		"o.deleted":         ctx.DefaultQuery("deleted", "0"),
		"o.currency":        ctx.Query("currency"),
	}

	var whereClauses []string
	var args []interface{}
	argCounter := 1

	// Basic filters
	for key, value := range filters {
		if value != "" && value != nil {
			whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", key, argCounter))
			args = append(args, value)
			argCounter++
		}
	}

	// Numeric range filters
	numericRanges := map[string]struct {
		min string
		max string
	}{
		"cost_per_km": {ctx.Query("min_cost_per_km"), ctx.Query("max_cost_per_km")},
		"offer_price": {ctx.Query("min_offer_price"), ctx.Query("max_offer_price")},
		"total_price": {ctx.Query("min_total_price"), ctx.Query("max_total_price")},
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

	// Date range filters
	validityStart := ctx.Query("validity_start")
	validityEnd := ctx.Query("validity_end")
	deliveryStart := ctx.Query("delivery_start")
	deliveryEnd := ctx.Query("delivery_end")

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

	// Search functionality
	searchTerm := ctx.Query("search")
	if searchTerm != "" {
		searchClause := fmt.Sprintf(`(
			o.note ILIKE $%d OR 
			o.meta ILIKE $%d OR
			o.from_address ILIKE $%d OR
			o.to_address ILIKE $%d OR
			o.from_country ILIKE $%d OR
			o.to_country ILIKE $%d OR
			o.from_region ILIKE $%d OR
			o.to_region ILIKE $%d OR
			o.sender_contact ILIKE $%d OR
			o.recipient_contact ILIKE $%d OR
			o.deliver_contact ILIKE $%d OR
			(d.first_name ILIKE $%d OR d.last_name ILIKE $%d OR d.patronymic_name ILIKE $%d OR 
			 d.phone ILIKE $%d OR d.email ILIKE $%d OR d.meta ILIKE $%d) OR
			(v.numberplate ILIKE $%d OR v.trailer_numberplate ILIKE $%d OR v.meta ILIKE $%d) OR
			(vt.numberplate ILIKE $%d OR vt.trailer_numberplate ILIKE $%d OR vt.meta ILIKE $%d)
		)`, argCounter, argCounter+1, argCounter+2, argCounter+3, argCounter+4, argCounter+5, argCounter+6, argCounter+7,
			argCounter+8, argCounter+9, argCounter+10, argCounter+11, argCounter+12, argCounter+13, argCounter+14, argCounter+15,
			argCounter+16, argCounter+17, argCounter+18, argCounter+19, argCounter+20, argCounter+21, argCounter+22)

		whereClauses = append(whereClauses, searchClause)
		searchPattern := "%" + searchTerm + "%"
		// Add 23 search parameters (11 offer fields + 6 driver fields + 6 vehicle/trailer fields)
		for i := 0; i < 23; i++ {
			args = append(args, searchPattern)
		}
		argCounter += 23
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT 
			o.*,
			COUNT(*) OVER() as total_count,
			to_json(c) as company_json,
			to_json(d) as driver_json,
			to_json(v) as vehicle_json,
			to_json(vt) as trailer_json,
			to_json(cr) as cargo_json,
			COALESCE((SELECT COUNT(*) FROM tbl_offer_response tor WHERE tor.offer_id = o.id AND tor.deleted = 0), 0) as response_count
		FROM tbl_offer o
		LEFT JOIN tbl_company c ON o.company_id = c.id
		LEFT JOIN tbl_driver d ON o.driver_id = d.id AND d.active = 1 AND d.deleted = 0
		LEFT JOIN tbl_vehicle v ON o.vehicle_id = v.id AND v.active = 1 AND v.deleted = 0
		LEFT JOIN tbl_vehicle vt ON o.trailer_id = vt.id AND vt.active = 1 AND vt.deleted = 0
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
		offer.Company = &dto.CompanyBasic{}
		offer.AssignedDriver = &dto.DriverShort{}
		offer.AssignedVehicle = &dto.VehicleBasic{}
		offer.AssignedTrailer = &dto.VehicleBasic{}
		offer.Cargo = &dto.CargoMain{}

		var companyJSON, driverJSON, vehicleJSON, trailerJSON, cargoJSON []byte

		err = rows.Scan(
			&offer.ID, &offer.UUID, &offer.UserID, &offer.CompanyID,
			&offer.ExecCompanyID, &offer.DriverID, &offer.VehicleID, &offer.TrailerID, &offer.VehicleTypeID,
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
			&offer.Deleted, &offer.OfferPrice, &offer.TotalPrice, &totalCount,
			&companyJSON, &driverJSON, &vehicleJSON, &trailerJSON, &cargoJSON,
			&offer.ResponseCount,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Scan error", err.Error()))
			return
		}

		// JSON unmarshaling logic remains the same
		if companyJSON != nil {
			if err = json.Unmarshal(companyJSON, &offer.Company); err != nil {
				ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("JSON unmarshal error (company)", err.Error()))
				return
			}
		} else {
			offer.Company = nil
		}

		if driverJSON != nil {
			if err = json.Unmarshal(driverJSON, &offer.AssignedDriver); err != nil {
				ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("JSON unmarshal error (driver)", err.Error()))
				return
			}
		} else {
			offer.AssignedDriver = nil
		}

		if vehicleJSON != nil {
			if err = json.Unmarshal(vehicleJSON, &offer.AssignedVehicle); err != nil {
				ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("JSON unmarshal error (vehicle)", err.Error()))
				return
			}
		} else {
			offer.AssignedVehicle = nil
		}

		if trailerJSON != nil {
			if err = json.Unmarshal(trailerJSON, &offer.AssignedTrailer); err != nil {
				ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("JSON unmarshal error (trailer)", err.Error()))
				return
			}
		} else {
			offer.AssignedTrailer = nil
		}

		if cargoJSON != nil {
			if err := json.Unmarshal(cargoJSON, &offer.Cargo); err != nil {
				ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("JSON unmarshal error (cargo)", err.Error()))
				return
			}
		} else {
			offer.Cargo = nil
		}

		offers = append(offers, offer)
	}

	if err = rows.Err(); err != nil {
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
		"trailer_id":      ctx.Query("trailer_id"),
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
	}

	if offer.OfferPrice == 0.0 {
		offer.OfferPrice = offer.CostPerKm * float64(offer.Distance)
	}
	if offer.TotalPrice == 0.0 {
		discountAmount := offer.OfferPrice * float64(offer.Discount) / 100
		offer.TotalPrice = offer.OfferPrice - discountAmount + offer.TaxPrice
	}

	var id int
	err := db.DB.QueryRow(
		context.Background(),
		queries.CreateOffer,
		offer.UserID, offer.CompanyID, offer.DriverID, offer.VehicleID, offer.TrailerID, offer.CargoID, offer.CostPerKm, offer.Currency,
		offer.FromCountryID, offer.FromCityID, offer.ToCountryID, offer.ToCityID,
		offer.FromCountry, offer.FromRegion, offer.ToCountry, offer.ToRegion, offer.FromAddress, offer.ToAddress,
		offer.SenderContact, offer.RecipientContact, offer.DeliverContact, offer.ValidityStart, offer.ValidityEnd,
		offer.DeliveryStart, offer.DeliveryEnd, offer.Note, offer.Tax, offer.TaxPrice, offer.Trade, offer.Discount,
		offer.PaymentMethod, offer.Meta, offer.Meta2, offer.Meta3, offer.OfferRole, offer.ExecCompanyID,
		offer.VehicleTypeID, offer.PackagingTypeID, offer.Distance, offer.MapURL, offer.PaymentTerm,
		offer.OfferPrice, offer.TotalPrice,
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
	companyID := ctx.MustGet("companyID").(int)
	isAdminOrSystem := (role == "admin" || role == "system")

	if !isAdminOrSystem {
		var ownershipCheck struct {
			CompanyID     *int `db:"company_id"`
			ExecCompanyID *int `db:"exec_company_id"`
		}

		checkQuery := `
            SELECT company_id, exec_company_id
            FROM tbl_offer
            WHERE id = $1 AND deleted = 0`

		err := db.DB.QueryRow(context.Background(), checkQuery, offerID).Scan(
			&ownershipCheck.CompanyID,
			&ownershipCheck.ExecCompanyID,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Offer not found", ""))
				return
			}
			ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error checking offer ownership", err.Error()))
			return
		}

		isOwner := false
		if ownershipCheck.CompanyID != nil && *ownershipCheck.CompanyID == companyID {
			isOwner = true
		}
		if ownershipCheck.ExecCompanyID != nil && *ownershipCheck.ExecCompanyID == companyID {
			isOwner = true
		}

		if !isOwner {
			ctx.JSON(http.StatusForbidden, utils.FormatErrorResponse("You don't have permission to update this offer", ""))
			return
		}

		offer.CompanyID = nil
		offer.OfferState = nil
		stmt += ` AND deleted = 0`
	}

	//
	//if offer.DriverID != nil {
	// var currentOffer struct {
	//    OfferRole     string `db:"offer_role"`
	//    CompanyID     int    `db:"company_id"`
	//    ExecCompanyID int    `db:"exec_company_id"`
	// }
	//
	// checkQuery := `
	//        SELECT offer_role, company_id, exec_company_id
	//        FROM tbl_offer
	//        WHERE id = $1 AND deleted = 0`
	//
	// err := db.DB.QueryRow(context.Background(), checkQuery, offerID).Scan(
	//    &currentOffer.OfferRole,
	//    &currentOffer.CompanyID,
	//    &currentOffer.ExecCompanyID,
	// )
	//
	// if err != nil {
	//    if err == sql.ErrNoRows {
	//       ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Offer not found", ""))
	//       return
	//    }
	//    ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error checking offer permissions", err.Error()))
	//    return
	// }
	//
	// canModifyDriver := false
	//
	// if isAdminOrSystem {
	//    canModifyDriver = true
	// } else if currentOffer.OfferRole == "carrier" && currentOffer.CompanyID == companyID {
	//    canModifyDriver = true
	// } else if currentOffer.OfferRole == "sender" && currentOffer.ExecCompanyID == companyID {
	//    canModifyDriver = true
	// }
	//
	// if !canModifyDriver {
	//    ctx.JSON(http.StatusForbidden, utils.FormatErrorResponse("You don't have permission to modify driver for this offer", ""))
	//    return
	// }
	//}

	if isAdminOrSystem {
		stmt = strings.Replace(stmt, "WHERE id = $1 AND company_id = $2", "WHERE id = $1", 1)
	} else {
		stmt = strings.Replace(stmt, "WHERE id = $1 AND company_id = $2", "WHERE id = $1 AND (company_id = $2 OR exec_company_id = $2)", 1)
	}

	stmt += ` RETURNING id;`

	var updatedID int
	err := db.DB.QueryRow(
		context.Background(),
		stmt,
		offerID,
		offer.CompanyID,
		offer.ExecCompanyID,
		offer.DriverID,
		offer.VehicleID,
		offer.VehicleTypeID,
		offer.CargoID,
		offer.PackagingTypeID,
		offer.OfferState,
		offer.OfferRole,
		offer.CostPerKm,
		offer.Currency,
		offer.FromCountryID,
		offer.FromCityID,
		offer.ToCountryID,
		offer.ToCityID,
		offer.Distance,
		offer.FromCountry,
		offer.FromRegion,
		offer.ToCountry,
		offer.ToRegion,
		offer.FromAddress,
		offer.ToAddress,
		offer.MapURL,
		offer.SenderContact,
		offer.RecipientContact,
		offer.DeliverContact,
		offer.ViewCount,
		offer.ValidityStart,
		offer.ValidityEnd,
		offer.DeliveryStart,
		offer.DeliveryEnd,
		offer.Note,
		offer.Tax,
		offer.TaxPrice,
		offer.Trade,
		offer.Discount,
		offer.PaymentMethod,
		offer.PaymentTerm,
		offer.Meta,
		offer.Meta2,
		offer.Meta3,
		offer.Featured,
		offer.Partner,
		offer.Active,
		offer.Deleted,
		offer.TrailerID,
		offer.OfferPrice,
		offer.TotalPrice,
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
				'id',c.id,
				'uuid',c.uuid,
				'user_id',c.user_id,
				'role',c.role,
				'role_id',c.role_id,
				'plan',c.plan,
				'plan_active',c.plan_active,
				'company_name',c.company_name,
				'first_name',c.first_name,
				'last_name',c.last_name,
				'patronymic_name',c.patronymic_name,
				'about',c.about,
				'phone',c.phone,
				'phone2',c.phone2,
				'phone3',c.phone3,
				'email',c.email,
				'email2',c.email2,
				'email3',c.email3,
				'meta',c.meta,
				'meta2',c.meta2,
				'meta3',c.meta3,
				'address',c.address,
				'country',c.country,
				'country_id',c.country_id,
				'city_id',c.city_id,
				'image_url',c.image_url,
				'verified',c.verified,
				'entity',c.entity,
				'featured',c.featured,
				'rating',c.rating,
				'partner',c.partner,
				'view_count',c.view_count,
				'successful_ops',c.successful_ops
            ) as company,
            json_build_object(
                'id',ec.id,
				'uuid',ec.uuid,
				'user_id',ec.user_id,
				'role',ec.role,
				'role_id',ec.role_id,
				'plan',ec.plan,
				'plan_active',ec.plan_active,
				'company_name',ec.company_name,
				'first_name',ec.first_name,
				'last_name',ec.last_name,
				'patronymic_name',ec.patronymic_name,
				'about',ec.about,
				'phone',ec.phone,
				'phone2',ec.phone2,
				'phone3',ec.phone3,
				'email',ec.email,
				'email2',ec.email2,
				'email3',ec.email3,
				'meta',ec.meta,
				'meta2',ec.meta2,
				'meta3',ec.meta3,
				'address',ec.address,
				'country',ec.country,
				'country_id',ec.country_id,
				'city_id',ec.city_id,
				'image_url',ec.image_url,
				'verified',ec.verified,
				'entity',ec.entity,
				'featured',ec.featured,
				'rating',ec.rating,
				'partner',ec.partner,
				'view_count',ec.view_count,
				'successful_ops',ec.successful_ops
            ) as exec_company,
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
				'available',d.available,
				'block_reason',d.block_reason,
				'active',d.active,
				'deleted',d.deleted
            ) as assigned_driver,
            json_build_object(
                'id', v.id,
				'uuid', v.uuid,
				'company_id', v.company_id,
				'vehicle_type_id', v.vehicle_type_id,
				'vehicle_brand_id', v.vehicle_brand_id,
				'vehicle_model_id', v.vehicle_model_id,
				'year_of_issue', v.year_of_issue,
				'mileage', v.mileage,
				'numberplate', v.numberplate,
				'trailer_numberplate', v.trailer_numberplate,
				'gps', v.gps,
				'photo1_url', v.photo1_url,
				'photo2_url', v.photo2_url,
				'photo3_url', v.photo3_url,
				'docs1_url', v.docs1_url,
				'docs2_url', v.docs2_url,
				'docs3_url', v.docs3_url,
				'view_count', v.view_count,
				'meta', v.meta,
				'meta2', v.meta2,
				'meta3', v.meta3,
				'available', v.available,
				'active', v.active,
				'deleted', v.deleted
            ) as vehicle,
			json_build_object(
                'id', vtr.id,
				'uuid', vtr.uuid,
				'company_id', vtr.company_id,
				'vehicle_type_id', vtr.vehicle_type_id,
				'vehicle_brand_id', vtr.vehicle_brand_id,
				'vehicle_model_id', vtr.vehicle_model_id,
				'year_of_issue', vtr.year_of_issue,
				'mileage', vtr.mileage,
				'numberplate', vtr.numberplate,
				'trailer_numberplate', vtr.trailer_numberplate,
				'gps', vtr.gps,
				'photo1_url', vtr.photo1_url,
				'photo2_url', vtr.photo2_url,
				'photo3_url', vtr.photo3_url,
				'docs1_url', vtr.docs1_url,
				'docs2_url', vtr.docs2_url,
				'docs3_url', vtr.docs3_url,
				'view_count', vtr.view_count,
				'meta', vtr.meta,
				'meta2', vtr.meta2,
				'meta3', vtr.meta3,
				'available', vtr.available,
				'active', vtr.active,
				'deleted', vtr.deleted
            ) as trailer,
            json_build_object(
                'id', vt.id,
				'title_en', vt.title_en,
				'desc_en', vt.desc_en,
				'title_ru', vt.title_ru,
				'desc_ru', vt.desc_ru,
				'title_tk', vt.title_tk,
				'desc_tk', vt.desc_tk,
				'title_de', vt.title_de,
				'desc_de', vt.desc_de,
				'title_ar', vt.title_ar,
				'desc_ar', vt.desc_ar,
				'title_es', vt.title_es,
				'desc_es', vt.desc_es,
				'title_fr', vt.title_fr,
				'desc_fr', vt.desc_fr,
				'title_zh', vt.title_zh,
				'desc_zh', vt.desc_zh,
				'title_ja', vt.title_ja,
				'desc_ja', vt.desc_ja,
				'deleted', vt.deleted
            ) as vehicle_type,
            json_build_object(
                'id', cg.id,
				'company_id', cg.company_id,
				'name', cg.name,
				'description', cg.description,
				'info', cg.info,
				'qty', cg.qty,
				'weight', cg.weight,
				'weight_type', cg.weight_type,
				'meta', cg.meta,
				'meta2', cg.meta2,
				'meta3', cg.meta3,
				'vehicle_type_id', cg.vehicle_type_id,
				'packaging_type_id', cg.packaging_type_id,
				'gps', cg.gps,
				'photo1_url', cg.photo1_url,
				'photo2_url', cg.photo2_url,
				'photo3_url', cg.photo3_url,
				'docs1_url', cg.docs1_url,
				'docs2_url', cg.docs2_url,
				'docs3_url', cg.docs3_url,
				'note', cg.note,
				'active', cg.active,
				'deleted', cg.deleted
            ) as cargo,
            json_build_object(
                'id', pt.id,
				'name_ru', pt.name_ru,
				'name_en', pt.name_en,
				'name_tk', pt.name_tk,
				'category_ru', pt.category_ru,
				'category_en', pt.category_en,
				'category_tk', pt.category_tk,
				'material', pt.material,
				'dimensions', pt.dimensions,
				'weight', pt.weight,
				'description_ru', pt.description_ru,
				'description_en', pt.description_en,
				'description_tk', pt.description_tk,
				'active', pt.active,
				'deleted', pt.deleted
            ) as packaging_type
        FROM tbl_offer o
        LEFT JOIN tbl_company c ON o.company_id = c.id
        LEFT JOIN tbl_company ec ON o.exec_company_id = ec.id
        LEFT JOIN tbl_driver d ON o.driver_id = d.id
        LEFT JOIN tbl_vehicle v ON o.vehicle_id = v.id
        LEFT JOIN tbl_vehicle vtr ON o.trailer_id = vtr.id
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

	if offer.CompanyID == ctx.MustGet("companyID") {
		const responseQuery = `
        SELECT 
            or_tbl.*,
            json_build_object(
                'id', c.id,
                'company_name', c.company_name,
                'first_name', c.first_name,
                'last_name', c.last_name
            ) as company,
            json_build_object(
                'id', tc.id,
                'company_name', tc.company_name,
                'first_name', tc.first_name,
                'last_name', tc.last_name
            ) as to_company
        FROM tbl_offer_response or_tbl
        LEFT JOIN tbl_company c ON or_tbl.company_id = c.id
        LEFT JOIN tbl_company tc ON or_tbl.to_company_id = tc.id
        WHERE or_tbl.offer_id = $1 AND or_tbl.deleted = 0
        ORDER BY or_tbl.created_at DESC`

		var offerResponses []dto.OfferResponseDetails
		err = pgxscan.Select(context.Background(), db.DB, &offerResponses, responseQuery, ctx.Param("id"))
		if err != nil && err != pgx.ErrNoRows {
			ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error fetching offer responses", err.Error()))
			return
		}

		offer.OfferResponses = offerResponses
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
		"offer_price": true, "total_price": true, "view_count": true,
		"validity_start": true, "validity_end": true, "delivery_start": true,
		"delivery_end": true, "created_at": true, "updated_at": true,
		"tax": true, "discount": true,
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
                    WHEN o2.deleted = 0 
                    AND o2.validity_end >= CURRENT_DATE 
                    THEN o2.id 
                    END) as offers_count
            FROM tbl_company c
            LEFT JOIN tbl_driver d ON d.company_id = c.id AND d.deleted = 0
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
				'uuid', v.uuid,
				'company_id', v.company_id,
				'vehicle_type_id', v.vehicle_type_id,
				'vehicle_brand_id', v.vehicle_brand_id,
				'vehicle_model_id', v.vehicle_model_id,
				'year_of_issue', v.year_of_issue,
				'mileage', v.mileage,
				'numberplate', v.numberplate,
				'trailer_numberplate', v.trailer_numberplate,
				'gps', v.gps,
				'photo1_url', v.photo1_url,
				'photo2_url', v.photo2_url,
				'photo3_url', v.photo3_url,
				'docs1_url', v.docs1_url,
				'docs2_url', v.docs2_url,
				'docs3_url', v.docs3_url,
				'view_count', v.view_count,
				'meta', v.meta,
				'meta2', v.meta2,
				'meta3', v.meta3,
				'available', v.available,
				'active', v.active,
				'deleted', v.deleted
            ) as vehicle,
			-- Trailer fields
            json_build_object(
				'id', vtr.id,
				'uuid', vtr.uuid,
				'company_id', vtr.company_id,
				'vehicle_type_id', vtr.vehicle_type_id,
				'vehicle_brand_id', vtr.vehicle_brand_id,
				'vehicle_model_id', vtr.vehicle_model_id,
				'year_of_issue', vtr.year_of_issue,
				'mileage', vtr.mileage,
				'numberplate', vtr.numberplate,
				'trailer_numberplate', vtr.trailer_numberplate,
				'gps', vtr.gps,
				'photo1_url', vtr.photo1_url,
				'photo2_url', vtr.photo2_url,
				'photo3_url', vtr.photo3_url,
				'docs1_url', vtr.docs1_url,
				'docs2_url', vtr.docs2_url,
				'docs3_url', vtr.docs3_url,
				'view_count', vtr.view_count,
				'meta', vtr.meta,
				'meta2', vtr.meta2,
				'meta3', vtr.meta3,
				'available', vtr.available,
				'active', vtr.active,
				'deleted', vtr.deleted
            ) as trailer,
            
            -- Vehicle type fields
            json_build_object(
                'id', vt.id,
				'title_en', vt.title_en,
				'desc_en', vt.desc_en,
				'title_ru', vt.title_ru,
				'desc_ru', vt.desc_ru,
				'title_tk', vt.title_tk,
				'desc_tk', vt.desc_tk,
				'title_de', vt.title_de,
				'desc_de', vt.desc_de,
				'title_ar', vt.title_ar,
				'desc_ar', vt.desc_ar,
				'title_es', vt.title_es,
				'desc_es', vt.desc_es,
				'title_fr', vt.title_fr,
				'desc_fr', vt.desc_fr,
				'title_zh', vt.title_zh,
				'desc_zh', vt.desc_zh,
				'title_ja', vt.title_ja,
				'desc_ja', vt.desc_ja,
				'deleted', vt.deleted
            ) as vehicle_type,
            
            -- Cargo fields
            json_build_object(
                'id', cg.id,
				'company_id', cg.company_id,
				'name', cg.name,
				'description', cg.description,
				'info', cg.info,
				'qty', cg.qty,
				'weight', cg.weight,
				'weight_type', cg.weight_type,
				'meta', cg.meta,
				'meta2', cg.meta2,
				'meta3', cg.meta3,
				'vehicle_type_id', cg.vehicle_type_id,
				'packaging_type_id', cg.packaging_type_id,
				'gps', cg.gps,
				'photo1_url', cg.photo1_url,
				'photo2_url', cg.photo2_url,
				'photo3_url', cg.photo3_url,
				'docs1_url', cg.docs1_url,
				'docs2_url', cg.docs2_url,
				'docs3_url', cg.docs3_url,
				'note', cg.note,
				'active', cg.active,
				'deleted', cg.deleted
            ) as cargo,
            
            -- Packaging type fields
            json_build_object(
                'id', pt.id,
				'name_ru', pt.name_ru,
				'name_en', pt.name_en,
				'name_tk', pt.name_tk,
				'category_ru', pt.category_ru,
				'category_en', pt.category_en,
				'category_tk', pt.category_tk,
				'material', pt.material,
				'dimensions', pt.dimensions,
				'weight', pt.weight,
				'description_ru', pt.description_ru,
				'description_en', pt.description_en,
				'description_tk', pt.description_tk,
				'active', pt.active,
				'deleted', pt.deleted
            ) as packaging_type
            
        FROM tbl_offer o
        LEFT JOIN tbl_company c ON o.company_id = c.id
        LEFT JOIN company_stats cs1 ON c.id = cs1.company_id
        LEFT JOIN tbl_company ec ON o.exec_company_id = ec.id
        LEFT JOIN company_stats cs2 ON ec.id = cs2.company_id
        LEFT JOIN tbl_driver d ON o.driver_id = d.id AND d.active = 1 AND d.deleted = 0
        LEFT JOIN tbl_vehicle v ON o.vehicle_id = v.id AND v.active = 1 AND v.deleted = 0
        LEFT JOIN tbl_vehicle vtr ON o.trailer_id = vtr.id AND vtr.active = 1 AND vtr.deleted = 0
        LEFT JOIN tbl_vehicle_type vt ON o.vehicle_type_id = vt.id
        LEFT JOIN tbl_cargo cg ON o.cargo_id = cg.id
        LEFT JOIN tbl_packaging_type pt ON o.packaging_type_id = pt.id
        `

	var whereClauses []string
	var args []interface{}
	argCounter := 1

	filters := map[string]string{
		"company_id":        ctx.Query("company_id"),
		"exec_company_id":   ctx.Query("exec_company_id"),
		"driver_id":         ctx.Query("driver_id"),
		"vehicle_id":        ctx.Query("vehicle_id"),
		"trailer_id":        ctx.Query("trailer_id"),
		"vehicle_type_id":   ctx.Query("vehicle_type_id"),
		"cargo_id":          ctx.Query("cargo_id"),
		"packaging_type_id": ctx.Query("packaging_type_id"),
		"offer_state":       ctx.Query("offer_state"),
		"offer_role":        ctx.Query("offer_role"),
		"currency":          ctx.Query("currency"),
		"payment_method":    ctx.Query("payment_method"),
		"active":            ctx.Query("active"),
		"deleted":           ctx.DefaultQuery("deleted", "0"),
		"featured":          ctx.Query("featured"),
		"partner":           ctx.Query("partner"),
		"from_country_id":   ctx.Query("from_country_id"),
		"from_city_id":      ctx.Query("from_city_id"),
		"to_country_id":     ctx.Query("to_country_id"),
		"to_city_id":        ctx.Query("to_city_id"),
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
		"cost_per_km": {ctx.Query("min_cost_per_km"), ctx.Query("max_cost_per_km")},
		"offer_price": {ctx.Query("min_offer_price"), ctx.Query("max_offer_price")},
		"total_price": {ctx.Query("min_total_price"), ctx.Query("max_total_price")},
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

	tripAddedCheck := ctx.Query("trip_added_check")
	if tripAddedCheck == "1" {
		excludeOffersQuery := `
        SELECT DISTINCT offer_id 
        FROM tbl_offer_trip 
        WHERE deleted = 0 AND status IN ('active', 'enabled')`

		var excludeOfferIds []int
		err := pgxscan.Select(
			context.Background(),
			db.DB,
			&excludeOfferIds,
			excludeOffersQuery,
		)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError,
				utils.FormatErrorResponse("Couldn't retrieve trip offers", err.Error()))
			return
		}

		if len(excludeOfferIds) > 0 {
			excludeIds := make([]string, len(excludeOfferIds))
			for i, id := range excludeOfferIds {
				excludeIds[i] = strconv.Itoa(id)
			}
			whereClauses = append(whereClauses, fmt.Sprintf("o.id NOT IN (%s)", strings.Join(excludeIds, ",")))
		}
	}

	// Global search
	searchTerm := ctx.Query("search")
	if searchTerm != "" {
		searchClause := fmt.Sprintf(`(
			o.note ILIKE $%d OR 
			o.meta ILIKE $%d OR
			o.from_address ILIKE $%d OR
			o.to_address ILIKE $%d OR
			o.from_country ILIKE $%d OR
			o.to_country ILIKE $%d OR
			o.from_region ILIKE $%d OR
			o.to_region ILIKE $%d OR
			o.sender_contact ILIKE $%d OR
			o.recipient_contact ILIKE $%d OR
			o.deliver_contact ILIKE $%d OR
			(d.first_name ILIKE $%d OR d.last_name ILIKE $%d OR d.patronymic_name ILIKE $%d OR 
			 d.phone ILIKE $%d OR d.email ILIKE $%d OR d.meta ILIKE $%d) OR
			(v.numberplate ILIKE $%d OR v.trailer_numberplate ILIKE $%d OR v.meta ILIKE $%d) OR
			(vtr.numberplate ILIKE $%d OR vtr.trailer_numberplate ILIKE $%d OR vtr.meta ILIKE $%d)
		)`, argCounter, argCounter+1, argCounter+2, argCounter+3, argCounter+4, argCounter+5, argCounter+6, argCounter+7,
			argCounter+8, argCounter+9, argCounter+10, argCounter+11, argCounter+12, argCounter+13, argCounter+14, argCounter+15,
			argCounter+16, argCounter+17, argCounter+18, argCounter+19, argCounter+20, argCounter+21, argCounter+22)

		whereClauses = append(whereClauses, searchClause)
		searchPattern := "%" + searchTerm + "%"
		// !! 23 search parameters (11 offer fields + 6 driver fields + 6 vehicle/trailer fields)
		for i := 0; i < 23; i++ {
			args = append(args, searchPattern)
		}
		argCounter += 23
	}

	searchFromLocation := ctx.Query("from_location")
	if searchFromLocation != "" {
		searchFromLocationClause := fmt.Sprintf(`(
        o.from_country ILIKE $%d OR
        o.from_region ILIKE $%d OR 
        o.from_address ILIKE $%d
    )`, argCounter, argCounter+1, argCounter+2)
		whereClauses = append(whereClauses, searchFromLocationClause)
		args = append(args, "%"+searchFromLocation+"%", "%"+searchFromLocation+"%", "%"+searchFromLocation+"%")
		argCounter += 3
	}

	searchToLocation := ctx.Query("to_location")
	if searchToLocation != "" {
		searchToLocationClause := fmt.Sprintf(`(
        o.to_country ILIKE $%d OR
        o.to_region ILIKE $%d OR 
        o.to_address ILIKE $%d
    )`, argCounter, argCounter+1, argCounter+2)
		whereClauses = append(whereClauses, searchToLocationClause)
		args = append(args, "%"+searchToLocation+"%", "%"+searchToLocation+"%", "%"+searchToLocation+"%")
		argCounter += 3
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
