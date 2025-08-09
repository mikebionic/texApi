package services

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"texApi/database"
	"texApi/internal/dto"
	"texApi/pkg/utils"
)

func GetFilteredClaims(ctx *gin.Context) {
	var filter dto.ClaimFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid filter parameters", err.Error()))
		return
	}

	claims, total, err := GetClaimsList(filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to retrieve claims", err.Error()))
		return
	}

	response := utils.PaginatedResponse{
		Total:   total,
		Page:    filter.Page,
		PerPage: filter.PerPage,
		Data:    claims,
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Claims retrieved successfully", response))
}

func NewClaim(ctx *gin.Context) {
	var req dto.NewClaimRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request data", err.Error()))
		return
	}

	claim, err := CreateClaimInternal(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to create claim", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.FormatResponse("Claim created successfully", claim))
}

func UpdateClaim(ctx *gin.Context) {
	var req dto.UpdateClaimRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request data", err.Error()))
		return
	}

	claim, err := UpdateClaimInternal(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to update claim", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Claim updated successfully", claim))
}

func DeleteClaim(ctx *gin.Context) {
	var req dto.DeleteClaimRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request data", err.Error()))
		return
	}

	if err := DeleteClaimInternal(req.ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to delete claim", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Claim deleted successfully", nil))
}

func GetClaimsList(filter dto.ClaimFilter) ([]dto.Claim, int, error) {
	var claims []dto.Claim
	var total int

	whereClauses := []string{"deleted = 0"}
	args := []interface{}{}
	argIndex := 1

	if filter.ClaimType != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("claim_type = $%d", argIndex))
		args = append(args, *filter.ClaimType)
		argIndex++
	}
	if filter.ClaimStatus != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("claim_status = $%d", argIndex))
		args = append(args, *filter.ClaimStatus)
		argIndex++
	}
	if filter.UrgencyLevel != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("urgency_level = $%d", argIndex))
		args = append(args, *filter.UrgencyLevel)
		argIndex++
	}
	if filter.UserID != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("user_id = $%d", argIndex))
		args = append(args, *filter.UserID)
		argIndex++
	}
	if filter.CompanyID != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("company_id = $%d", argIndex))
		args = append(args, *filter.CompanyID)
		argIndex++
	}
	if filter.Email != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("email = $%d", argIndex))
		args = append(args, *filter.Email)
		argIndex++
	}
	if filter.Active != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("active = $%d", argIndex))
		args = append(args, *filter.Active)
		argIndex++
	}
	if filter.Search != nil {
		searchClause := fmt.Sprintf("(subject ILIKE $%d OR description ILIKE $%d OR name ILIKE $%d OR email ILIKE $%d)",
			argIndex, argIndex+1, argIndex+2, argIndex+3)
		whereClauses = append(whereClauses, searchClause)
		searchTerm := "%" + *filter.Search + "%"
		args = append(args, searchTerm, searchTerm, searchTerm, searchTerm)
		argIndex += 4
	}

	whereClause := strings.Join(whereClauses, " AND ")

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM tbl_claim WHERE %s", whereClause)
	err := database.DB.QueryRow(context.Background(), countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	query := fmt.Sprintf(`
		SELECT id, uuid, user_id, company_id, name, email, phone, address, company_name,
		       subject, description, additional_details, response_title, response_description,
		       claim_type, claim_status, urgency_level, meta, meta2, meta3,
		       created_at, updated_at, active, deleted
		FROM tbl_claim
		WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`,
		whereClause, argIndex, argIndex+1)

	args = append(args, filter.PerPage, (filter.Page-1)*filter.PerPage)
	err = pgxscan.Select(context.Background(), database.DB, &claims, query, args...)
	return claims, total, err
}

func CreateClaimInternal(req dto.NewClaimRequest) (dto.Claim, error) {
	var claim dto.Claim
	query := `
		INSERT INTO tbl_claim (
			user_id, company_id, name, email, phone, address, company_name,
			subject, description, additional_details, claim_type, urgency_level,
			meta, meta2, meta3
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		) RETURNING id, uuid, user_id, company_id, name, email, phone, address, company_name,
		           subject, description, additional_details, response_title, response_description,
		           claim_type, claim_status, urgency_level, meta, meta2, meta3,
		           created_at, updated_at, active, deleted`

	err := pgxscan.Get(context.Background(), database.DB, &claim, query,
		req.UserID, req.CompanyID, req.Name, req.Email, req.Phone, req.Address,
		req.CompanyName, req.Subject, req.Description, req.AdditionalDetails,
		req.ClaimType, req.UrgencyLevel, req.Meta, req.Meta2, req.Meta3)
	return claim, err
}

func UpdateClaimInternal(req dto.UpdateClaimRequest) (dto.Claim, error) {
	var claim dto.Claim
	setParts := []string{"updated_at = CURRENT_TIMESTAMP"}
	args := []interface{}{}
	argIndex := 1

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
	if req.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *req.Name)
		argIndex++
	}
	if req.Email != nil {
		setParts = append(setParts, fmt.Sprintf("email = $%d", argIndex))
		args = append(args, *req.Email)
		argIndex++
	}
	if req.Phone != nil {
		setParts = append(setParts, fmt.Sprintf("phone = $%d", argIndex))
		args = append(args, *req.Phone)
		argIndex++
	}
	if req.Address != nil {
		setParts = append(setParts, fmt.Sprintf("address = $%d", argIndex))
		args = append(args, *req.Address)
		argIndex++
	}
	if req.CompanyName != nil {
		setParts = append(setParts, fmt.Sprintf("company_name = $%d", argIndex))
		args = append(args, *req.CompanyName)
		argIndex++
	}
	if req.Subject != nil {
		setParts = append(setParts, fmt.Sprintf("subject = $%d", argIndex))
		args = append(args, *req.Subject)
		argIndex++
	}
	if req.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, *req.Description)
		argIndex++
	}
	if req.AdditionalDetails != nil {
		setParts = append(setParts, fmt.Sprintf("additional_details = $%d", argIndex))
		args = append(args, *req.AdditionalDetails)
		argIndex++
	}
	if req.ResponseTitle != nil {
		setParts = append(setParts, fmt.Sprintf("response_title = $%d", argIndex))
		args = append(args, *req.ResponseTitle)
		argIndex++
	}
	if req.ResponseDescription != nil {
		setParts = append(setParts, fmt.Sprintf("response_description = $%d", argIndex))
		args = append(args, *req.ResponseDescription)
		argIndex++
	}
	if req.ClaimType != nil {
		setParts = append(setParts, fmt.Sprintf("claim_type = $%d", argIndex))
		args = append(args, *req.ClaimType)
		argIndex++
	}
	if req.ClaimStatus != nil {
		setParts = append(setParts, fmt.Sprintf("claim_status = $%d", argIndex))
		args = append(args, *req.ClaimStatus)
		argIndex++
	}
	if req.UrgencyLevel != nil {
		setParts = append(setParts, fmt.Sprintf("urgency_level = $%d", argIndex))
		args = append(args, *req.UrgencyLevel)
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

	args = append(args, req.ID)
	setClause := strings.Join(setParts, ", ")
	query := fmt.Sprintf(`
		UPDATE tbl_claim 
		SET %s 
		WHERE id = $%d AND deleted = 0 
		RETURNING id, uuid, user_id, company_id, name, email, phone, address, company_name,
		          subject, description, additional_details, response_title, response_description,
		          claim_type, claim_status, urgency_level, meta, meta2, meta3,
		          created_at, updated_at, active, deleted`, setClause, argIndex)

	err := pgxscan.Get(context.Background(), database.DB, &claim, query, args...)
	return claim, err
}

func DeleteClaimInternal(id int) error {
	query := `UPDATE tbl_claim SET deleted = 1, updated_at = CURRENT_TIMESTAMP WHERE id = $1 AND deleted = 0`
	_, err := database.DB.Exec(context.Background(), query, id)
	return err
}
