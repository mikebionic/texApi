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
                'validity_start', o.validity_start,
                'validity_end', o.validity_end
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
	var offerResponse dto.OfferResponseCreate

	if err := ctx.ShouldBindJSON(&offerResponse); err != nil {
		ctx.JSON(http.StatusBadRequest,
			utils.FormatErrorResponse("Invalid request body", err.Error()))
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
	err := db.DB.QueryRow(
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
func UpdateOfferResponse(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var offerResponse dto.OfferResponseUpdate

	userID := ctx.MustGet("id").(int)
	role := ctx.MustGet("role").(string)

	if err := ctx.ShouldBindJSON(&offerResponse); err != nil {
		ctx.JSON(http.StatusBadRequest,
			utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	// If not admin, verify company ownership
	if role != "admin" {
		var count int
		err := db.DB.QueryRow(
			context.Background(),
			`SELECT COUNT(*) 
             FROM tbl_offer_response ofr
             JOIN tbl_company c ON ofr.company_id = c.id
             WHERE ofr.id = $1 AND c.user_id = $2 AND ofr.deleted = 0`,
			id, userID,
		).Scan(&count)

		if err != nil || count == 0 {
			ctx.JSON(http.StatusForbidden,
				utils.FormatErrorResponse("Permission denied", ""))
			return
		}

		// Reset admin-only fields for non-admin users
		offerResponse.Active = nil
		offerResponse.Deleted = nil
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
            deleted = COALESCE($11, deleted),
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $12 AND deleted = 0
        RETURNING id
    `

	var updatedID int
	err := db.DB.QueryRow(
		context.Background(),
		query,
		offerResponse.State, offerResponse.BidPrice,
		offerResponse.Title, offerResponse.Note,
		offerResponse.Reason, offerResponse.Meta,
		offerResponse.Meta2, offerResponse.Meta3,
		offerResponse.Value, offerResponse.Rating,
		offerResponse.Active, offerResponse.Deleted,
		id,
	).Scan(&updatedID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			utils.FormatErrorResponse("Error updating offer response", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully updated!", gin.H{
		"id": updatedID,
	}))
}

func DeleteOfferResponse(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	userID := ctx.MustGet("id").(int)
	role := ctx.MustGet("role").(string)

	// If not admin, verify company ownership
	if role != "admin" {
		var count int
		err := db.DB.QueryRow(
			context.Background(),
			`SELECT COUNT(*) 
             FROM tbl_offer_response ofr
             JOIN tbl_company c ON ofr.company_id = c.id
             WHERE ofr.id = $1 AND c.user_id = $2 AND ofr.deleted = 0`,
			id, userID,
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
        WHERE id = $1 AND deleted = 0
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
