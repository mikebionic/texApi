package services

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	db "texApi/database"
	"texApi/internal/dto"
	"texApi/internal/queries"
	"texApi/pkg/utils"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
)

func Subscribe(ctx *gin.Context) {
	var req dto.NewsletterSubscribe

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	if req.IPAddress == "" {
		req.IPAddress = ctx.ClientIP()
	}
	if req.UserAgent == "" {
		req.UserAgent = ctx.GetHeader("User-Agent")
	}
	if req.ReferrerURL == "" {
		req.ReferrerURL = ctx.GetHeader("Referer")
	}

	if req.Frequency == "" {
		req.Frequency = "other"
	}

	var existingID int
	err := db.DB.QueryRow(context.Background(), queries.CheckEmailExists, req.Email).Scan(&existingID)
	if err == nil {
		ctx.JSON(http.StatusConflict, utils.FormatErrorResponse("Email already subscribed", ""))
		return
	}

	var id int
	err = db.DB.QueryRow(
		context.Background(),
		queries.CreateNewsletter,
		req.Email, "active", req.FirstName, req.LastName, req.Frequency,
		req.IPAddress, req.UserAgent, req.ReferrerURL, req.Meta, req.Meta2, req.Meta3,
	).Scan(&id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error subscribing to newsletter", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.FormatResponse("Successfully subscribed to newsletter!", gin.H{"id": id}))
}

func GetNewsletterList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
	offset := (page - 1) * perPage

	validOrderColumns := map[string]bool{
		"id": true, "email": true, "status": true, "subscribed_at": true,
		"created_at": true, "updated_at": true, "first_name": true, "last_name": true,
	}

	orderBy := ctx.DefaultQuery("order_by", "id")
	if !validOrderColumns[orderBy] {
		orderBy = "id"
	}
	orderDir := strings.ToUpper(ctx.DefaultQuery("order_dir", "DESC"))
	if orderDir != "ASC" && orderDir != "DESC" {
		orderDir = "DESC"
	}

	baseQuery := queries.GetNewsletterList
	var whereClauses []string
	var args []interface{}
	argCounter := 1

	filters := map[string]string{
		"status":    ctx.Query("status"),
		"frequency": ctx.Query("frequency"),
		"active":    ctx.Query("active"),
		"deleted":   ctx.Query("deleted"),
	}

	for key, value := range filters {
		if value != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("n.%s = $%d", key, argCounter))
			args = append(args, value)
			argCounter++
		}
	}

	searchTerm := ctx.Query("search")
	if searchTerm != "" {
		searchClause := fmt.Sprintf(`(
			n.email ILIKE $%d OR 
			n.first_name ILIKE $%d OR 
			n.last_name ILIKE $%d
		)`, argCounter, argCounter, argCounter)
		whereClauses = append(whereClauses, searchClause)
		args = append(args, "%"+searchTerm+"%")
		argCounter++
	}

	if ctx.Query("deleted") == "" {
		whereClauses = append(whereClauses, "n.deleted = 0")
	}

	query := baseQuery
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	query += fmt.Sprintf(" ORDER BY n.%s %s LIMIT $%d OFFSET $%d",
		orderBy, orderDir, argCounter, argCounter+1)
	args = append(args, perPage, offset)

	var newsletters []dto.Newsletter
	err := pgxscan.Select(
		context.Background(),
		db.DB,
		&newsletters,
		query,
		args...,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Couldn't retrieve data", err.Error()))
		return
	}

	var totalCount int
	if len(newsletters) > 0 {
		totalCount = newsletters[0].TotalCount
	}

	response := utils.PaginatedResponse{
		Total:   totalCount,
		Page:    page,
		PerPage: perPage,
		Data:    newsletters,
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Newsletter list", response))
}

func UpdateNewsletter(ctx *gin.Context) {
	id := ctx.Param("id")
	var newsletter dto.NewsletterUpdate

	if err := ctx.ShouldBindJSON(&newsletter); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	if newsletter.Status != nil {
		validStatuses := map[string]bool{
			"active": true, "unsubscribed": true, "bounced": true, "pending": true,
		}
		if !validStatuses[*newsletter.Status] {
			ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid status", "Valid statuses: active, unsubscribed, bounced, pending"))
			return
		}
	}

	if newsletter.Frequency != nil {
		validFrequencies := map[string]bool{
			"daily": true, "weekly": true, "monthly": true, "other": true,
		}
		if !validFrequencies[*newsletter.Frequency] {
			ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid frequency", "Valid frequencies: daily, weekly, monthly, other"))
			return
		}
	}

	result, err := db.DB.Exec(
		context.Background(),
		queries.UpdateNewsletter,
		id, newsletter.Email, newsletter.Status, newsletter.FirstName, newsletter.LastName,
		newsletter.Frequency, newsletter.IPAddress, newsletter.UserAgent, newsletter.ReferrerURL,
		newsletter.Meta, newsletter.Meta2, newsletter.Meta3, newsletter.Active, newsletter.Deleted,
	)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			ctx.JSON(http.StatusConflict, utils.FormatErrorResponse("Email already exists", ""))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error updating newsletter", err.Error()))
		return
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Newsletter not found or no changes were made", ""))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully updated newsletter!", gin.H{"id": id}))
}

func DeleteNewsletter(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := db.DB.Exec(
		context.Background(),
		queries.DeleteNewsletter,
		id,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error deleting newsletter", err.Error()))
		return
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Newsletter not found", ""))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully deleted newsletter!", gin.H{"id": id}))
}
