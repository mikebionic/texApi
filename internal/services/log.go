package services

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"texApi/internal/queries"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"

	db "texApi/database"
	"texApi/internal/dto"
	"texApi/pkg/utils"
)

func GetVerifyRequestList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
	offset := (page - 1) * perPage

	userID := ctx.Query("user_id")
	companyID := ctx.Query("company_id")
	status := ctx.Query("status")
	createdDateFrom := ctx.Query("created_date_from")
	createdDateTo := ctx.Query("created_date_to")
	createdDate := ctx.Query("created_date")
	orderBy := ctx.DefaultQuery("order_by", "id")
	orderDir := ctx.DefaultQuery("order_dir", "ASC")

	baseQuery := `
                SELECT vr.*, u.username, p.first_name, p.last_name, p.company_name, p.role, count(*) OVER() as total_count 
                FROM tbl_verify_request vr
                LEFT JOIN tbl_user u ON vr.user_id = u.id
                LEFT JOIN tbl_company p ON vr.company_id = p.id
                WHERE vr.deleted = 0
        `

	validOrderColumns := map[string]bool{
		"id": true, "user_id": true, "company_id": true, "status": true, "created_at": true,
	}

	conditions := []string{}
	args := []interface{}{}
	paramCount := 1

	if userID != "" {
		conditions = append(conditions, fmt.Sprintf("vr.user_id = $%d", paramCount))
		args = append(args, userID)
		paramCount++
	}

	if companyID != "" {
		conditions = append(conditions, fmt.Sprintf("vr.company_id = $%d", paramCount))
		args = append(args, companyID)
		paramCount++
	}

	if status != "" {
		conditions = append(conditions, fmt.Sprintf("vr.status = $%d", paramCount))
		args = append(args, status)
		paramCount++
	}

	if createdDateFrom != "" {
		conditions = append(conditions, fmt.Sprintf("vr.created_at >= $%d", paramCount))
		args = append(args, createdDateFrom)
		paramCount++
	}

	if createdDateTo != "" {
		conditions = append(conditions, fmt.Sprintf("vr.created_at <= $%d", paramCount))
		args = append(args, createdDateTo)
		paramCount++
	}

	if createdDate != "" {
		conditions = append(conditions, fmt.Sprintf("vr.created_at::DATE = $%d", paramCount))
		args = append(args, createdDate)
		paramCount++
	}

	// Append conditions to query
	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}

	// Add ordering and pagination
	if validOrderColumns[orderBy] {
		baseQuery += fmt.Sprintf(" ORDER BY %s %s", orderBy, orderDir)
	} else {
		baseQuery += " ORDER BY id ASC"
	}
	baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramCount, paramCount+1)
	args = append(args, perPage, offset)

	var requests []dto.VerifyRequestDetails
	err := pgxscan.Select(context.Background(), db.DB, &requests, baseQuery, args...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Database error", err.Error()))
		return
	}

	totalCount := 0
	if len(requests) > 0 {
		totalCount = requests[0].TotalCount
	}

	response := utils.PaginatedResponse{
		Total:   totalCount,
		Page:    page,
		PerPage: perPage,
		Data:    requests,
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Verification request list", response))
}

func GetVerifyRequest(ctx *gin.Context) {
	id := ctx.Param("id")

	query := `
		SELECT * 
		FROM tbl_verify_request 
		WHERE id = $1 AND deleted = 0
	`

	var request dto.VerifyRequestDetails
	err := pgxscan.Get(context.Background(), db.DB, &request, query, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Verification request not found", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Verification request details", request))
}

// CreateVerifyRequest creates a new verification request USER's ROUTE!
func CreateVerifyRequest(ctx *gin.Context) {
	userID := ctx.MustGet("id").(int)
	companyID := ctx.MustGet("companyID").(int)

	// Check if there's an existing pending request for this company
	checkQuery := `
		SELECT COUNT(*) 
		FROM tbl_verify_request 
		WHERE user_id=$1 AND company_id = $2 AND status = 'pending' AND deleted = 0
	`
	var count int
	err := db.DB.QueryRow(context.Background(), checkQuery, userID, companyID).Scan(&count)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Database error", err.Error()))
		return
	}

	if count > 0 {
		ctx.JSON(http.StatusConflict, utils.FormatErrorResponse("A pending verification request already exists for this company", ""))
		return
	}

	query := `
		INSERT INTO tbl_verify_request (
			user_id, company_id
		) VALUES (
			$1, $2
		) RETURNING id
	`

	var requestID int
	err = db.DB.QueryRow(
		context.Background(),
		query,
		userID, companyID,
	).Scan(&requestID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating verification request", err.Error()))
		return
	}

	// Log the action
	CreateUserLog(context.Background(), dto.UserLogCreate{
		UserID:    userID,
		CompanyID: companyID,
		Role:      ctx.GetString("role"),
		Action:    "create_verify_request",
		Details:   fmt.Sprintf("Created verification request ID: %d", requestID),
	})

	ctx.JSON(http.StatusCreated, utils.FormatResponse("Successfully created verification request!", gin.H{"id": requestID}))
}

// UpdateVerifyRequest updates a verification request (admin only)
func UpdateVerifyRequest(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var request dto.VerifyRequestUpdate

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	// First, get the current request to get the company_id
	getQuery := `
		SELECT company_id, user_id
		FROM tbl_verify_request
		WHERE id = $1 AND deleted = 0
	`
	var companyID, userID int
	err := db.DB.QueryRow(context.Background(), getQuery, id).Scan(&companyID, &userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Verification request not found", err.Error()))
		return
	}

	// Update the request status
	updateQuery := `
		UPDATE tbl_verify_request
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND deleted = 0
	`

	commandTag, err := db.DB.Exec(
		context.Background(),
		updateQuery,
		request.Status,
		id,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error updating verification request", err.Error()))
		return
	}

	rowsAffected := commandTag.RowsAffected()
	if rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Verification request not found or already processed", ""))
		return
	}

	currentUserRole := ctx.MustGet("role").(string)
	currentUserID := ctx.MustGet("id").(int)

	// If approved, update the company verification status
	if request.Status != nil && *request.Status == "approved" {
		_, err := queries.UpdateCompanyVerification(context.Background(), db.DB, companyID, 1)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error updating company verification status", err.Error()))
			return
		}

		// Log the approval
		CreateUserLog(context.Background(), dto.UserLogCreate{
			UserID:    currentUserID,
			CompanyID: companyID,
			Role:      currentUserRole,
			Action:    "approve_verify_request",
			Details:   fmt.Sprintf("Approved verification request ID: %d", id),
		})
	} else if request.Status != nil && *request.Status == "declined" {
		// Log the decline
		CreateUserLog(context.Background(), dto.UserLogCreate{
			UserID:    currentUserID,
			CompanyID: companyID,
			Role:      currentUserRole,
			Action:    "decline_verify_request",
			Details:   fmt.Sprintf("Declined verification request ID: %d", id),
		})
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully updated verification request!", gin.H{"id": id}))
}

// DeleteVerifyRequest soft deletes a verification request (admin only)
func DeleteVerifyRequest(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	// Get the request details before deletion
	getQuery := `
		SELECT company_id, user_id
		FROM tbl_verify_request
		WHERE id = $1 AND deleted = 0
	`
	var companyID, userID int
	err := db.DB.QueryRow(context.Background(), getQuery, id).Scan(&companyID, &userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Verification request not found", err.Error()))
		return
	}

	// Soft delete
	query := `
		UPDATE tbl_verify_request 
		SET deleted = 1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND deleted = 0
	`

	commandTag, err := db.DB.Exec(
		context.Background(),
		query,
		id,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error deleting verification request", err.Error()))
		return
	}

	rowsAffected := commandTag.RowsAffected()
	if rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Verification request not found or already deleted", ""))
		return
	}

	// Log the deletion
	currentUserRole := ctx.MustGet("role").(string)
	currentUserID := ctx.MustGet("id").(int)

	CreateUserLog(context.Background(), dto.UserLogCreate{
		UserID:    currentUserID,
		CompanyID: companyID,
		Role:      currentUserRole,
		Action:    "delete_verify_request",
		Details:   fmt.Sprintf("Deleted verification request ID: %d", id),
	})

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully deleted verification request!", gin.H{"id": id}))
}

// GetPlanMovesList retrieves a paginated list of plan moves
func GetPlanMovesList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
	offset := (page - 1) * perPage

	userID := ctx.Query("user_id")
	companyID := ctx.Query("company_id")
	status := ctx.Query("status")
	createdDateFrom := ctx.Query("created_date_from")
	createdDateTo := ctx.Query("created_date_to")
	createdDate := ctx.Query("created_date")
	validUntilFrom := ctx.Query("valid_until_from")
	validUntilTo := ctx.Query("valid_until_to")
	orderBy := ctx.DefaultQuery("order_by", "id")
	orderDir := ctx.DefaultQuery("order_dir", "ASC")

	baseQuery := `
                SELECT pm.*, u.username, p.first_name, p.last_name, p.company_name, p.role, count(*) OVER() as total_count 
                FROM tbl_plan_moves pm
                LEFT JOIN tbl_user u ON pm.user_id = u.id
                LEFT JOIN tbl_company p ON pm.company_id = p.id
                WHERE pm.deleted = 0
        `

	validOrderColumns := map[string]bool{
		"id": true, "user_id": true, "company_id": true, "status": true, "valid_until": true, "created_at": true,
	}

	conditions := []string{}
	args := []interface{}{}
	paramCount := 1

	if userID != "" {
		conditions = append(conditions, fmt.Sprintf("pm.user_id = $%d", paramCount))
		args = append(args, userID)
		paramCount++
	}

	if companyID != "" {
		conditions = append(conditions, fmt.Sprintf("pm.company_id = $%d", paramCount))
		args = append(args, companyID)
		paramCount++
	}

	if status != "" {
		conditions = append(conditions, fmt.Sprintf("pm.status = $%d", paramCount))
		args = append(args, status)
		paramCount++
	}

	if createdDateFrom != "" {
		conditions = append(conditions, fmt.Sprintf("pm.created_at >= $%d", paramCount))
		args = append(args, createdDateFrom)
		paramCount++
	}

	if createdDateTo != "" {
		conditions = append(conditions, fmt.Sprintf("pm.created_at <= $%d", paramCount))
		args = append(args, createdDateTo)
		paramCount++
	}

	if createdDate != "" {
		conditions = append(conditions, fmt.Sprintf("pm.created_at::DATE = $%d", paramCount))
		args = append(args, createdDate)
		paramCount++
	}

	if validUntilFrom != "" {
		conditions = append(conditions, fmt.Sprintf("pm.valid_until >= $%d", paramCount))
		args = append(args, validUntilFrom)
		paramCount++
	}

	if validUntilTo != "" {
		conditions = append(conditions, fmt.Sprintf("pm.valid_until <= $%d", paramCount))
		args = append(args, validUntilTo)
		paramCount++
	}

	// Append conditions to query
	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}

	// Add ordering and pagination
	if validOrderColumns[orderBy] {
		baseQuery += fmt.Sprintf(" ORDER BY %s %s", orderBy, orderDir)
	} else {
		baseQuery += " ORDER BY id ASC"
	}
	baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramCount, paramCount+1)
	args = append(args, perPage, offset)

	var planMoves []dto.PlanMoveDetails
	err := pgxscan.Select(context.Background(), db.DB, &planMoves, baseQuery, args...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Database error", err.Error()))
		return
	}

	totalCount := 0
	if len(planMoves) > 0 {
		totalCount = planMoves[0].TotalCount
	}

	response := utils.PaginatedResponse{
		Total:   totalCount,
		Page:    page,
		PerPage: perPage,
		Data:    planMoves,
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Plan moves list", response))
}

// GetPlanMove retrieves a single plan move by ID
func GetPlanMove(ctx *gin.Context) {
	id := ctx.Param("id")

	query := `
		SELECT * 
		FROM tbl_plan_moves 
		WHERE id = $1 AND deleted = 0
	`

	var planMove dto.PlanMoveDetails
	err := pgxscan.Get(context.Background(), db.DB, &planMove, query, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Plan move not found", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Plan move details", planMove))
}

// CreatePlanMove creates a new plan move
func CreatePlanMove(ctx *gin.Context) {
	var planMove dto.PlanMoveCreate

	if err := ctx.ShouldBindJSON(&planMove); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	// Check if there's an existing pending plan move for this company
	checkQuery := `
		SELECT COUNT(*) 
		FROM tbl_plan_moves 
		WHERE company_id = $1 AND status = 'pending' AND deleted = 0
	`
	var count int
	err := db.DB.QueryRow(context.Background(), checkQuery, planMove.CompanyID).Scan(&count)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Database error", err.Error()))
		return
	}

	if count > 0 {
		ctx.JSON(http.StatusConflict, utils.FormatErrorResponse("A pending plan move already exists for this company", ""))
		return
	}

	// Default status is 'pending' if not provided
	status := planMove.Status
	if status == "" {
		status = "pending"
	}

	query := `
		INSERT INTO tbl_plan_moves (
			user_id, company_id, status
		) VALUES (
			$1, $2, $3::status_type_t
		) RETURNING id
	`

	var planMoveID int
	err = db.DB.QueryRow(
		context.Background(),
		query,
		planMove.UserID, planMove.CompanyID, status,
	).Scan(&planMoveID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error creating plan move", err.Error()))
		return
	}

	// Log the action
	currentUserRole := ctx.GetString("role")
	currentUserID := ctx.MustGet("id").(int)

	CreateUserLog(context.Background(), dto.UserLogCreate{
		UserID:    currentUserID,
		CompanyID: planMove.CompanyID,
		Role:      currentUserRole,
		Action:    "create_plan_move",
		Details:   fmt.Sprintf("Created plan move ID: %d", planMoveID),
	})

	ctx.JSON(http.StatusCreated, utils.FormatResponse("Successfully created plan move!", gin.H{"id": planMoveID}))
}

// UpdatePlanMove updates a plan move (admin only)
func UpdatePlanMove(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var planMove dto.PlanMoveUpdate

	if err := ctx.ShouldBindJSON(&planMove); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request body", err.Error()))
		return
	}

	// First, get the current plan move details
	getQuery := `
		SELECT company_id, user_id, status
		FROM tbl_plan_moves
		WHERE id = $1 AND deleted = 0
	`
	var companyID, userID int
	var currentStatus string
	err := db.DB.QueryRow(context.Background(), getQuery, id).Scan(&companyID, &userID, &currentStatus)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Plan move not found", err.Error()))
		return
	}

	tx, err := db.DB.Begin(context.Background())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Transaction error", err.Error()))
		return
	}
	defer tx.Rollback(context.Background())

	// If we're changing to approved status, update the valid_until date
	if planMove.Status != nil && *planMove.Status == "approved" && currentStatus != "approved" {
		validUntil, err := queries.ExtendPlanValidity(context.Background(), db.DB, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error updating plan validity", err.Error()))
			return
		}

		// Update the company's plan_active status
		_, err = queries.UpdateCompanyPlanActive(context.Background(), db.DB, companyID, 1)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error updating company plan status", err.Error()))
			return
		}

		// Set the valid_until time for the update
		planMove.ValidUntil = &validUntil
	}

	// Build the update query based on what fields are provided
	updateQuery := `
		UPDATE tbl_plan_moves
		SET updated_at = CURRENT_TIMESTAMP
	`
	args := []interface{}{}
	paramCount := 0

	if planMove.Status != nil {
		paramCount++
		updateQuery += fmt.Sprintf(", status = $%d::status_type_t", paramCount)
		args = append(args, *planMove.Status)
	}

	if planMove.ValidUntil != nil {
		paramCount++
		updateQuery += fmt.Sprintf(", valid_until = $%d", paramCount)
		args = append(args, *planMove.ValidUntil)
	}

	paramCount++
	updateQuery += fmt.Sprintf(" WHERE id = $%d AND deleted = 0", paramCount)
	args = append(args, id)

	commandTag, err := db.DB.Exec(context.Background(), updateQuery, args...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error updating plan move", err.Error()))
		return
	}

	rowsAffected := commandTag.RowsAffected()
	if rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Plan move not found or already processed", ""))
		return
	}

	// Log the action
	currentUserRole := ctx.GetString("role")
	currentUserID := ctx.MustGet("id").(int)

	action := "update_plan_move"
	details := fmt.Sprintf("Updated plan move ID: %d", id)

	if planMove.Status != nil && *planMove.Status == "approved" {
		action = "approve_plan_move"
		details = fmt.Sprintf("Approved plan move ID: %d", id)
	} else if planMove.Status != nil && *planMove.Status == "declined" {
		action = "decline_plan_move"
		details = fmt.Sprintf("Declined plan move ID: %d", id)
	}

	CreateUserLog(context.Background(), dto.UserLogCreate{
		UserID:    currentUserID,
		CompanyID: companyID,
		Role:      currentUserRole,
		Action:    action,
		Details:   details,
	})

	if err := tx.Commit(context.Background()); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error committing transaction", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully updated plan move!", gin.H{"id": id}))
}

// DeletePlanMove soft deletes a plan move (admin only)
func DeletePlanMove(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	// Get the plan move details before deletion
	getQuery := `
		SELECT company_id, user_id
		FROM tbl_plan_moves
		WHERE id = $1 AND deleted = 0
	`
	var companyID, userID int
	err := db.DB.QueryRow(context.Background(), getQuery, id).Scan(&companyID, &userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Plan move not found", err.Error()))
		return
	}

	// Soft delete
	query := `
		UPDATE tbl_plan_moves 
		SET deleted = 1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND deleted = 0
	`

	commandTag, err := db.DB.Exec(
		context.Background(),
		query,
		id,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error deleting plan move", err.Error()))
		return
	}

	rowsAffected := commandTag.RowsAffected()
	if rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Plan move not found or already deleted", ""))
		return
	}

	// Log the deletion
	currentUserRole := ctx.GetString("role")
	currentUserID := ctx.MustGet("id").(int)

	CreateUserLog(context.Background(), dto.UserLogCreate{
		UserID:    currentUserID,
		CompanyID: companyID,
		Role:      currentUserRole,
		Action:    "delete_plan_move",
		Details:   fmt.Sprintf("Deleted plan move ID: %d", id),
	})

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully deleted plan move!", gin.H{"id": id}))
}

// CheckExpiredPlans checks for and handles expired plans
func CheckExpiredPlans(ctx *gin.Context) {
	// This function could be called by a scheduler or admin endpoint
	rowsAffected, err := queries.CheckExpiredPlans(context.Background(), db.DB)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error checking expired plans", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully checked expired plans", gin.H{"updated_companys": rowsAffected}))
}

// CreateUserLog creates a new user log entry
func CreateUserLog(ctx context.Context, logData dto.UserLogCreate) (int, error) {
	query := `
		INSERT INTO tbl_user_log (
			user_id, company_id, role, action, details
		) VALUES (
			$1, $2, $3::role_t, $4, $5
		) RETURNING id
	`

	var logID int
	err := db.DB.QueryRow(
		ctx,
		query,
		logData.UserID, logData.CompanyID, logData.Role, logData.Action, logData.Details,
	).Scan(&logID)

	if err != nil {
		return 0, err
	}

	return logID, nil
}

func GetUserLogsList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
	offset := (page - 1) * perPage

	userID := ctx.Query("user_id")
	companyID := ctx.Query("company_id")
	role := ctx.Query("role")
	action := ctx.Query("action")
	searchDetails := ctx.Query("search_details")
	createdDateFrom := ctx.Query("created_date_from")
	createdDateTo := ctx.Query("created_date_to")
	createdDate := ctx.Query("created_date")
	orderBy := ctx.DefaultQuery("order_by", "id")
	orderDir := ctx.DefaultQuery("order_dir", "DESC")

	baseQuery := `
		SELECT ul.*, u.username, p.first_name, p.last_name, p.company_name, count(*) OVER() as total_count
		FROM tbl_user_log ul
		LEFT JOIN tbl_user u ON ul.user_id = u.id
		LEFT JOIN tbl_company p ON ul.company_id = p.id
		WHERE ul.deleted = 0
	`

	validOrderColumns := map[string]bool{
		"id": true, "user_id": true, "company_id": true, "role": true, "action": true, "created_at": true,
	}

	conditions := []string{}
	args := []interface{}{}
	paramCount := 1

	if userID != "" {
		conditions = append(conditions, fmt.Sprintf("ul.user_id = $%d", paramCount))
		args = append(args, userID)
		paramCount++
	}

	if companyID != "" {
		conditions = append(conditions, fmt.Sprintf("ul.company_id = $%d", paramCount))
		args = append(args, companyID)
		paramCount++
	}

	if role != "" {
		conditions = append(conditions, fmt.Sprintf("ul.role = $%d", paramCount))
		args = append(args, role)
		paramCount++
	}

	if action != "" {
		conditions = append(conditions, fmt.Sprintf("ul.action = $%d", paramCount))
		args = append(args, action)
		paramCount++
	}

	if searchDetails != "" {
		conditions = append(conditions, fmt.Sprintf("ul.details LIKE $%d", paramCount))
		args = append(args, "%"+searchDetails+"%")
		paramCount++
	}

	if createdDateFrom != "" {
		conditions = append(conditions, fmt.Sprintf("ul.created_at >= $%d", paramCount))
		args = append(args, createdDateFrom)
		paramCount++
	}

	if createdDateTo != "" {
		conditions = append(conditions, fmt.Sprintf("ul.created_at <= $%d", paramCount))
		args = append(args, createdDateTo)
		paramCount++
	}

	if createdDate != "" {
		conditions = append(conditions, fmt.Sprintf("ul.created_at::DATE = $%d", paramCount))
		args = append(args, createdDate)
		paramCount++
	}

	// Append conditions to query
	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}

	// Add ordering and pagination
	if validOrderColumns[orderBy] {
		baseQuery += fmt.Sprintf(" ORDER BY %s %s", orderBy, orderDir)
	} else {
		baseQuery += " ORDER BY id DESC"
	}
	baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramCount, paramCount+1)
	args = append(args, perPage, offset)

	var logs []dto.UserLogDetails
	err := pgxscan.Select(context.Background(), db.DB, &logs, baseQuery, args...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Database error", err.Error()))
		return
	}

	// Get total count
	totalCount := 0
	if len(logs) > 0 {
		totalCount = logs[0].TotalCount
	}

	response := utils.PaginatedResponse{
		Total:   totalCount,
		Page:    page,
		PerPage: perPage,
		Data:    logs,
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("User logs list", response))
}

// GetUserLog retrieves a single user log by ID
func GetUserLog(ctx *gin.Context) {
	id := ctx.Param("id")

	query := `
		SELECT ul.*, u.username, p.first_name, p.last_name, p.company_name, count(*) OVER() as total_count
		FROM tbl_user_log ul
		LEFT JOIN tbl_user u ON ul.user_id = u.id
		LEFT JOIN tbl_company p ON ul.company_id = p.id
		WHERE ul.id = $1 AND ul.deleted = 0
	`

	var log dto.UserLogDetails
	err := pgxscan.Get(context.Background(), db.DB, &log, query, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("User log not found", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("User log details", log))
}

// DeleteUserLog soft deletes a user log (admin only)
func DeleteUserLog(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	// Soft delete
	query := `
            UPDATE tbl_user_log 
            SET deleted = 1, updated_at = CURRENT_TIMESTAMP
            WHERE id = $1 AND deleted = 0
    `

	commandTag, err := db.DB.Exec(
		context.Background(),
		query,
		id,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error deleting user log", err.Error()))
		return
	}

	rowsAffected := commandTag.RowsAffected()
	if rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("User log not found or already deleted", ""))
		return
	}

	//// Log the deletion action itself
	//currentUserRole := ctx.GetString("role")
	//currentUserID := ctx.MustGet("id").(int)
	//
	//CreateUserLog(context.Background(), dto.UserLogCreate{
	//	UserID:    currentUserID,
	//	CompanyID: 0, // No specific company for this action
	//	Role:      currentUserRole,
	//	Action:    "delete_user_log",
	//	Details:   fmt.Sprintf("Deleted user log ID: %d", id),
	//})

	ctx.JSON(http.StatusOK, utils.FormatResponse("Successfully deleted user log!", gin.H{"id": id}))
}
