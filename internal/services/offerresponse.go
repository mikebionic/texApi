package services

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	db "texApi/database"
	"texApi/internal/dto"
	"texApi/pkg/utils"
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
				'id', c.id,
				'uuid', c.uuid,
				'company_name', c.company_name,
				'first_name', c.first_name,
				'last_name', c.last_name,
				'phone', c.phone,
				'email', c.email,
				'address', c.address,
				'country', c.country,
				'image_url', c.image_url,
				'rating', c.rating,
				'partner', c.partner
			) as company,
			
			json_build_object(
				'id', tc.id,
				'uuid', tc.uuid,
				'company_name', tc.company_name,
				'first_name', tc.first_name,
				'last_name', tc.last_name,
				'phone', tc.phone,
				'email', tc.email,
				'address', tc.address,
				'country', tc.country,
				'image_url', tc.image_url,
				'rating', tc.rating,
				'partner', tc.partner
			) as to_company,
			
			json_build_object(
				'id', o.id,
				'uuid', o.uuid,
				'offer_state', o.offer_state,
				'offer_role', o.offer_role,
				'cost_per_km', o.cost_per_km,
				'currency', o.currency,
				'distance', o.distance,
				'from_country', o.from_country,
				'to_country', o.to_country,
				'from_address', o.from_address,
				'to_address', o.to_address,
				'validity_start', to_char(o.validity_start, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
				'validity_end', to_char(o.validity_end, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
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

	// Override company_id with authenticated company's ID and set initial state
	offerResponse.CompanyID = companyID
	offerResponse.State = "pending"

	// Verify that the offer exists and is active
	var offerExists int
	err := db.DB.QueryRow(
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

	ctx.JSON(http.StatusCreated, utils.FormatResponse("Successfully created!", gin.H{
		"id":   responseID,
		"uuid": responseUUID,
	}))
}

// Updated UpdateOfferResponse function
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

	// Start transaction for updating offer responses
	tx, err := db.DB.Begin(context.Background())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			utils.FormatErrorResponse("Transaction error", err.Error()))
		return
	}
	defer tx.Rollback(context.Background())

	// If accepting the offer, decline all other responses for the same offer
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

	// Update the current offer response
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

	// Commit the transaction
	if err := tx.Commit(context.Background()); err != nil {
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

	// If not admin, verify company ownership
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

	// Soft delete the offer response
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
