package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
	"texApi/config"
	db "texApi/database"
	"texApi/internal/dto"
	"texApi/internal/repo"
	"texApi/pkg/utils"
	"time"
)

func GetDetailedOfferResponseList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
	offset := (page - 1) * perPage

	validOrderColumns := map[string]bool{
		"id": true, "bid_price": true, "rating": true,
		"value": true, "created_at": true, "updated_at": true,
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
        SELECT 
			ofr.*,
			COUNT(*) OVER() as total_count,
			
			json_build_object(			
				'id',c.id,
				'uuid',c.uuid,
				'user_id',c.user_id,
				'role',c.role,
				'role_id',c.role_id,
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
				'featured',c.featured,
				'rating',c.rating,
				'partner',c.partner
			) as company,
			
			json_build_object(
				'id',tc.id,
				'uuid',tc.uuid,
				'user_id',tc.user_id,
				'role',tc.role,
				'role_id',tc.role_id,
				'company_name',tc.company_name,
				'first_name',tc.first_name,
				'last_name',tc.last_name,
				'patronymic_name',tc.patronymic_name,
				'about',tc.about,
				'phone',tc.phone,
				'phone2',tc.phone2,
				'phone3',tc.phone3,
				'email',tc.email,
				'email2',tc.email2,
				'email3',tc.email3,
				'meta',tc.meta,
				'meta2',tc.meta2,
				'meta3',tc.meta3,
				'address',tc.address,
				'country',tc.country,
				'country_id',tc.country_id,
				'city_id',tc.city_id,
				'image_url',tc.image_url,
				'featured',tc.featured,
				'rating',tc.rating,
				'partner',tc.partner
			) as to_company,
			
			json_build_object(
				'id', o.id,
				'uuid', o.uuid,
				'user_id', o.user_id,
				'company_id', o.company_id,
				'exec_company_id', o.exec_company_id,
				'driver_id', o.driver_id,
				'vehicle_id', o.vehicle_id,
				'trailer_id', o.trailer_id,
				'vehicle_type_id', o.vehicle_type_id,
				'cargo_id', o.cargo_id,
				'packaging_type_id', o.packaging_type_id,
				'offer_state', o.offer_state,
				'offer_role', o.offer_role,
				'cost_per_km', o.cost_per_km,
				'currency', o.currency,
				'from_country_id', o.from_country_id,
				'from_city_id', o.from_city_id,
				'to_country_id', o.to_country_id,
				'to_city_id', o.to_city_id,
				'distance', o.distance,
				'from_country', o.from_country,
				'from_region', o.from_region,
				'to_country', o.to_country,
				'to_region', o.to_region,
				'from_address', o.from_address,
				'to_address', o.to_address,
				'map_url', o.map_url,
				'sender_contact', o.sender_contact,
				'recipient_contact', o.recipient_contact,
				'deliver_contact', o.deliver_contact,
				'view_count', o.view_count,
				'validity_start', to_char(o.validity_start, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
				'validity_end', to_char(o.validity_end, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),				
				'delivery_start', to_char(o.delivery_start, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
				'delivery_end', to_char(o.delivery_end, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
				'note', o.note,
				'tax', o.tax,
				'tax_price', o.tax_price,
				'trade', o.trade,
				'discount', o.discount,
				'payment_method', o.payment_method,
				'payment_term', o.payment_term,
				'meta', o.meta,
				'meta2', o.meta2,
				'meta3', o.meta3,
				'featured', o.featured,
				'partner', o.partner
			) as offer
			
		FROM tbl_offer_response ofr
		LEFT JOIN tbl_company c ON ofr.company_id = c.id
		LEFT JOIN tbl_company tc ON ofr.to_company_id = tc.id
		LEFT JOIN tbl_offer o ON ofr.offer_id = o.id
    `

	var whereClauses []string
	var args []interface{}
	argCounter := 1

	role := ctx.MustGet("role").(string)
	if !(role == "admin" || role == "system") {
		whereClauses = append(whereClauses, "ofr.deleted = 0")
	}

	filters := map[string]string{
		"company_id":    ctx.Query("company_id"),
		"to_company_id": ctx.Query("to_company_id"),
		"offer_id":      ctx.Query("offer_id"),
		"state":         ctx.Query("state"),
	}

	for key, value := range filters {
		if value != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("ofr.%s = $%d", key, argCounter))
			args = append(args, value)
			argCounter++
		}
	}

	numericRanges := map[string]struct {
		min string
		max string
	}{
		"bid_price": {ctx.Query("min_bid_price"), ctx.Query("max_bid_price")},
		"value":     {ctx.Query("min_value"), ctx.Query("max_value")},
		"rating":    {ctx.Query("min_rating"), ctx.Query("max_rating")},
	}

	for field, ranges := range numericRanges {
		if ranges.min != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("ofr.%s >= $%d", field, argCounter))
			minVal, _ := strconv.ParseFloat(ranges.min, 64)
			args = append(args, minVal)
			argCounter++
		}
		if ranges.max != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("ofr.%s <= $%d", field, argCounter))
			maxVal, _ := strconv.ParseFloat(ranges.max, 64)
			args = append(args, maxVal)
			argCounter++
		}
	}

	searchTerm := ctx.Query("search")
	if searchTerm != "" {
		searchClause := fmt.Sprintf(`(
            ofr.title ILIKE $%d OR 
            ofr.note ILIKE $%d OR 
            ofr.reason ILIKE $%d
        )`, argCounter, argCounter, argCounter)
		whereClauses = append(whereClauses, searchClause)
		args = append(args, "%"+searchTerm+"%")
		argCounter++
	}

	query := baseQuery
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	query += fmt.Sprintf(" ORDER BY ofr.%s %s LIMIT $%d OFFSET $%d",
		orderBy, orderDir, argCounter, argCounter+1)
	args = append(args, perPage, offset)

	var responses []dto.OfferResponseDetails
	err := pgxscan.Select(context.Background(), db.DB, &responses, query, args...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			utils.FormatErrorResponse("Database error", err.Error()))
		return
	}

	totalCount := 0
	if len(responses) > 0 {
		totalCount = responses[0].TotalCount
	}

	response := utils.PaginatedResponse{
		Total:   totalCount,
		Page:    page,
		PerPage: perPage,
		Data:    responses,
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Offer response list", response))
}

func GetOfferResponse(ctx *gin.Context) {
	id := ctx.Param("id")

	query := `
        SELECT ofr.*, 
            json_build_object(
                'id', c.id,
                'company_name', c.company_name,
                'email', c.email,
                'phone', c.phone
            ) as company,
            json_build_object(
                'id', tc.id,
                'company_name', tc.company_name,
                'email', tc.email,
                'phone', tc.phone
            ) as to_company,
            json_build_object(
                'id', o.id,
                'offer_state', o.offer_state,
                'cost_per_km', o.cost_per_km,
                'currency', o.currency
            ) as offer
        FROM tbl_offer_response ofr
        LEFT JOIN tbl_company c ON ofr.company_id = c.id
        LEFT JOIN tbl_company tc ON ofr.to_company_id = tc.id
        LEFT JOIN tbl_offer o ON ofr.offer_id = o.id
        WHERE (ofr.id::TEXT = $1 OR ofr.uuid::TEXT = $1) AND ofr.deleted = 0
    `

	var response dto.OfferResponseDetails
	err := pgxscan.Get(context.Background(), db.DB, &response, query, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound,
			utils.FormatErrorResponse("Offer response not found", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Offer response details", response))
}

func CreateOfferResponse(ctx *gin.Context) {
	var offerResponse dto.OfferResponse
	companyID := ctx.MustGet("companyID").(int)

	if err := ctx.ShouldBindJSON(&offerResponse); err != nil {
		ctx.JSON(http.StatusBadRequest,
			utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	offerResponse.CompanyID = companyID
	offerResponse.State = "pending"

	company, err := repo.GetCompanyByID(offerResponse.ToCompanyID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error retrieving recipient company", err.Error()))
		return
	}

	var offerExists int
	err = db.DB.QueryRow(
		context.Background(),
		`SELECT COUNT(*) FROM tbl_offer WHERE id = $1 AND deleted = 0`,
		offerResponse.OfferID,
	).Scan(&offerExists)

	if err != nil || offerExists == 0 {
		ctx.JSON(http.StatusBadRequest,
			utils.FormatErrorResponse("Invalid offer ID", "Offer not found or inactive"))
		return
	}

	query := `
        INSERT INTO tbl_offer_response (
            company_id, offer_id, to_company_id, state,
            bid_price, title, note, reason,
            meta, meta2, meta3, value, rating
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
        ) RETURNING id, uuid
    `

	var responseID int
	var responseUUID string
	err = db.DB.QueryRow(
		context.Background(),
		query,
		offerResponse.CompanyID, offerResponse.OfferID,
		offerResponse.ToCompanyID, offerResponse.State,
		offerResponse.BidPrice, offerResponse.Title,
		offerResponse.Note, offerResponse.Reason,
		offerResponse.Meta, offerResponse.Meta2,
		offerResponse.Meta3, offerResponse.Value,
		offerResponse.Rating,
	).Scan(&responseID, &responseUUID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			utils.FormatErrorResponse("Error creating offer response", err.Error()))
		return
	}

	go sendOfferResponseNotification(company.UserID, fmt.Sprintf("New offer response: %s", utils.SafeString(offerResponse.Title)), offerResponse)

	ctx.JSON(http.StatusCreated, utils.FormatResponse("Successfully created!", gin.H{
		"id":   responseID,
		"uuid": responseUUID,
	}))
}

func sendOfferResponseNotification(userID int, content string, data interface{}) {
	payload := map[string]interface{}{
		"userID":  userID,
		"content": content,
		"extras":  data,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal notification payload: %s", err)
		return
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("http://localhost:%s/%s/ws-notification/", config.ENV.API_PORT, config.ENV.API_PREFIX),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		log.Printf("Failed to create notification request: %s", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(config.ENV.SYSTEM_HEADER, config.ENV.API_SECRET)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send notification request %s", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf(fmt.Sprintf("Notification API returned non-OK status: %d", resp.StatusCode), nil)
	}
}

// Accept or Decline Offer Response
func UpdateOfferResponse(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var offerResponse dto.OfferResponseUpdate
	companyID := ctx.MustGet("companyID").(int)

	if err := ctx.ShouldBindJSON(&offerResponse); err != nil {
		ctx.JSON(http.StatusBadRequest,
			utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	// Verify the user belongs to to_company_id
	var toCompanyID int
	err := db.DB.QueryRow(
		context.Background(),
		`SELECT to_company_id FROM tbl_offer_response WHERE id = $1 AND deleted = 0`,
		id,
	).Scan(&toCompanyID)

	if err != nil {
		ctx.JSON(http.StatusNotFound,
			utils.FormatErrorResponse("Offer response not found", err.Error()))
		return
	}

	if toCompanyID != companyID {
		ctx.JSON(http.StatusForbidden,
			utils.FormatErrorResponse("Permission denied", "Only the recipient company can update this response"))
		return
	}

	tx, err := db.DB.Begin(context.Background())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			utils.FormatErrorResponse("Transaction error", err.Error()))
		return
	}
	defer tx.Rollback(context.Background())

	if offerResponse.State != nil && *offerResponse.State == "accepted" {
		var offerID int
		err := tx.QueryRow(
			context.Background(),
			`SELECT offer_id FROM tbl_offer_response WHERE id = $1`,
			id,
		).Scan(&offerID)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError,
				utils.FormatErrorResponse("Error fetching offer ID", err.Error()))
			return
		}

		_, err = tx.Exec(
			context.Background(),
			`UPDATE tbl_offer_response 
             SET state = 'declined', updated_at = CURRENT_TIMESTAMP
             WHERE offer_id = $1 AND id != $2 AND state = 'pending'`,
			offerID, id,
		)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError,
				utils.FormatErrorResponse("Error updating other responses", err.Error()))
			return
		}
	}

	query := `
        UPDATE tbl_offer_response SET
            state = COALESCE($1, state),
            bid_price = COALESCE($2, bid_price),
            title = COALESCE($3, title),
            note = COALESCE($4, note),
            reason = COALESCE($5, reason),
            meta = COALESCE($6, meta),
            meta2 = COALESCE($7, meta2),
            meta3 = COALESCE($8, meta3),
            value = COALESCE($9, value),
            rating = COALESCE($10, rating),
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $11 AND deleted = 0
        RETURNING id
    `

	var updatedID int
	err = tx.QueryRow(
		context.Background(),
		query,
		offerResponse.State, offerResponse.BidPrice,
		offerResponse.Title, offerResponse.Note,
		offerResponse.Reason, offerResponse.Meta,
		offerResponse.Meta2, offerResponse.Meta3,
		offerResponse.Value, offerResponse.Rating,
		id,
	).Scan(&updatedID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			utils.FormatErrorResponse("Error updating offer response", err.Error()))
		return
	}

	if err = tx.Commit(context.Background()); err != nil {
		ctx.JSON(http.StatusInternalServerError,
			utils.FormatErrorResponse("Error committing transaction", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully updated!", gin.H{
		"id": updatedID,
	}))
}
func DeleteOfferResponse(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	companyID := ctx.MustGet("companyID").(int)
	role := ctx.MustGet("role").(string)

	if role != "admin" {
		var count int
		err := db.DB.QueryRow(
			context.Background(),
			`SELECT COUNT(*) 
             FROM tbl_offer_response ofr
             WHERE ofr.id = $1 AND ofr.company_id = $2 AND ofr.state = 'pending'`,
			id, companyID,
		).Scan(&count)

		if err != nil || count == 0 {
			ctx.JSON(http.StatusForbidden,
				utils.FormatErrorResponse("Permission denied", ""))
			return
		}
	}

	query := `
        UPDATE tbl_offer_response 
        SET deleted = 1, updated_at = CURRENT_TIMESTAMP
        WHERE id = $1
        RETURNING id
    `

	var deletedID int
	err := db.DB.QueryRow(context.Background(), query, id).Scan(&deletedID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			utils.FormatErrorResponse("Error deleting offer response", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully deleted!", gin.H{
		"id": deletedID,
	}))
}
