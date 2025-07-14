package services

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"texApi/database"
	"texApi/internal/dto"
	"texApi/pkg/utils"
)

func GetPlans(ctx *gin.Context) {
	var filter dto.PlanFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid filter parameters", err.Error()))
		return
	}

	plans, total, err := GetPlansList(filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to retrieve plans", err.Error()))
		return
	}

	response := utils.PaginatedResponse{
		Total:   total,
		Page:    filter.Page,
		PerPage: filter.PerPage,
		Data:    plans,
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Plans", response))
}

func GetPlanByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("uuid"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid plan ID", err.Error()))
		return
	}

	plan, err := GetPlanByIDInternal(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Plan not found", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Plan retrieved", plan))
}

func CreatePlan(ctx *gin.Context) {
	var req dto.CreatePlanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request data", err.Error()))
		return
	}

	plan, err := CreatePlanInternal(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to create plan", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.FormatResponse("Plan created successfully", plan))
}

func UpdatePlan(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("uuid"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid plan ID", err.Error()))
		return
	}

	var req dto.UpdatePlanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request data", err.Error()))
		return
	}

	plan, err := UpdatePlanInternal(id, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to update plan", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Plan updated successfully", plan))
}

func DeletePlan(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("uuid"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid plan ID", err.Error()))
		return
	}

	if err := DeletePlanInternal(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to delete plan", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Plan deleted successfully", nil))
}

func GetPlansList(filter dto.PlanFilter) ([]dto.Plan, int, error) {
	var plans []dto.Plan
	var total int

	whereClauses := []string{"deleted = 0"}
	args := []interface{}{}
	argIndex := 1

	if filter.Provider != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("provider = $%d", argIndex))
		args = append(args, *filter.Provider)
		argIndex++
	}
	if filter.Region != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("region = $%d", argIndex))
		args = append(args, *filter.Region)
		argIndex++
	}
	if filter.BillingCycle != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("billing_cycle = $%d", argIndex))
		args = append(args, *filter.BillingCycle)
		argIndex++
	}
	if filter.IsPopular != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("is_popular = $%d", argIndex))
		args = append(args, *filter.IsPopular)
		argIndex++
	}
	if filter.IsRecommended != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("is_recommended = $%d", argIndex))
		args = append(args, *filter.IsRecommended)
		argIndex++
	}
	if filter.Active != nil {
		if *filter.Active {
			whereClauses = append(whereClauses, "active = 1")
		} else {
			whereClauses = append(whereClauses, "active = 0")
		}
	}
	if filter.Search != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("(name ILIKE $%d OR code ILIKE $%d OR description ILIKE $%d)", argIndex, argIndex, argIndex))
		args = append(args, "%"+*filter.Search+"%")
		argIndex++
	}

	whereClause := strings.Join(whereClauses, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM tbl_plan WHERE %s", whereClause)
	err := database.DB.QueryRow(context.Background(), countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT id, uuid, name, code, provider, region, price_usd, price_local, 
		       local_currency, billing_cycle, load_posts_limit, load_posts_unlimited,
		       gps_tracking_level, gps_has_eta, rate_tools_level, rate_tools_features,
		       edocs_available, edocs_limit, edocs_has_archiving, support_level,
		       payment_guarantee, api_access, display_order, is_popular, is_recommended,
		       description, features_summary, available_from, available_until,
		       meta, meta2, meta3, created_at, updated_at
		FROM tbl_plan
		WHERE %s
		ORDER BY display_order ASC, name ASC
		LIMIT $%d OFFSET $%d`,
		whereClause, argIndex, argIndex+1)

	args = append(args, filter.PerPage, (filter.Page-1)*filter.PerPage)
	err = pgxscan.Select(context.Background(), database.DB, &plans, query, args...)
	return plans, total, err
}

func GetPlanByIDInternal(id uuid.UUID) (dto.Plan, error) {
	var plan dto.Plan
	query := `
		SELECT id, uuid, name, code, provider, region, price_usd, price_local, 
		       local_currency, billing_cycle, load_posts_limit, load_posts_unlimited,
		       gps_tracking_level, gps_has_eta, rate_tools_level, rate_tools_features,
		       edocs_available, edocs_limit, edocs_has_archiving, support_level,
		       payment_guarantee, api_access, display_order, is_popular, is_recommended,
		       description, features_summary, available_from, available_until,
		       meta, meta2, meta3, created_at, updated_at
		FROM tbl_plan
		WHERE uuid = $1 AND deleted = 0`
	err := pgxscan.Get(context.Background(), database.DB, &plan, query, id)
	return plan, err
}

func CreatePlanInternal(req dto.CreatePlanRequest) (dto.Plan, error) {
	var plan dto.Plan
	query := `
		INSERT INTO tbl_plan (
			name, code, provider, region, price_usd, price_local, local_currency,
			billing_cycle, load_posts_limit, load_posts_unlimited, gps_tracking_level,
			gps_has_eta, rate_tools_level, rate_tools_features, edocs_available,
			edocs_limit, edocs_has_archiving, support_level, payment_guarantee,
			api_access, display_order, is_popular, is_recommended, description,
			features_summary, available_from, available_until, meta, meta2, meta3
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
			$16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30
		) RETURNING id, uuid, name, code, provider, region, price_usd, price_local, 
		           local_currency, billing_cycle, load_posts_limit, load_posts_unlimited,
		           gps_tracking_level, gps_has_eta, rate_tools_level, rate_tools_features,
		           edocs_available, edocs_limit, edocs_has_archiving, support_level,
		           payment_guarantee, api_access, display_order, is_popular, is_recommended,
		           description, features_summary, available_from, available_until,
		           meta, meta2, meta3, created_at, updated_at`

	err := pgxscan.Get(context.Background(), database.DB, &plan, query,
		req.Name, req.Code, req.Provider, req.Region, req.PriceUSD, req.PriceLocal,
		req.LocalCurrency, req.BillingCycle, req.LoadPostsLimit, req.LoadPostsUnlimited,
		req.GPSTrackingLevel, req.GPSHasETA, req.RateToolsLevel, req.RateToolsFeatures,
		req.EdocsAvailable, req.EdocsLimit, req.EdocsHasArchiving, req.SupportLevel,
		req.PaymentGuarantee, req.APIAccess, req.DisplayOrder, req.IsPopular,
		req.IsRecommended, req.Description, req.FeaturesSummary, req.AvailableFrom,
		req.AvailableUntil, req.Meta, req.Meta2, req.Meta3)
	return plan, err
}

func UpdatePlanInternal(id uuid.UUID, req dto.UpdatePlanRequest) (dto.Plan, error) {
	var plan dto.Plan
	setParts := []string{"updated_at = NOW()"}
	args := []interface{}{}
	argIndex := 1

	if req.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *req.Name)
		argIndex++
	}
	if req.Code != nil {
		setParts = append(setParts, fmt.Sprintf("code = $%d", argIndex))
		args = append(args, *req.Code)
		argIndex++
	}
	if req.Provider != nil {
		setParts = append(setParts, fmt.Sprintf("provider = $%d", argIndex))
		args = append(args, *req.Provider)
		argIndex++
	}
	if req.Region != nil {
		setParts = append(setParts, fmt.Sprintf("region = $%d", argIndex))
		args = append(args, *req.Region)
		argIndex++
	}
	if req.PriceUSD != nil {
		setParts = append(setParts, fmt.Sprintf("price_usd = $%d", argIndex))
		args = append(args, *req.PriceUSD)
		argIndex++
	}
	if req.PriceLocal != nil {
		setParts = append(setParts, fmt.Sprintf("price_local = $%d", argIndex))
		args = append(args, *req.PriceLocal)
		argIndex++
	}
	if req.LocalCurrency != nil {
		setParts = append(setParts, fmt.Sprintf("local_currency = $%d", argIndex))
		args = append(args, *req.LocalCurrency)
		argIndex++
	}
	if req.BillingCycle != nil {
		setParts = append(setParts, fmt.Sprintf("billing_cycle = $%d", argIndex))
		args = append(args, *req.BillingCycle)
		argIndex++
	}
	if req.LoadPostsLimit != nil {
		setParts = append(setParts, fmt.Sprintf("load_posts_limit = $%d", argIndex))
		args = append(args, *req.LoadPostsLimit)
		argIndex++
	}
	if req.LoadPostsUnlimited != nil {
		setParts = append(setParts, fmt.Sprintf("load_posts_unlimited = $%d", argIndex))
		args = append(args, *req.LoadPostsUnlimited)
		argIndex++
	}
	if req.GPSTrackingLevel != nil {
		setParts = append(setParts, fmt.Sprintf("gps_tracking_level = $%d", argIndex))
		args = append(args, *req.GPSTrackingLevel)
		argIndex++
	}
	if req.GPSHasETA != nil {
		setParts = append(setParts, fmt.Sprintf("gps_has_eta = $%d", argIndex))
		args = append(args, *req.GPSHasETA)
		argIndex++
	}
	if req.RateToolsLevel != nil {
		setParts = append(setParts, fmt.Sprintf("rate_tools_level = $%d", argIndex))
		args = append(args, *req.RateToolsLevel)
		argIndex++
	}
	if req.RateToolsFeatures != nil {
		setParts = append(setParts, fmt.Sprintf("rate_tools_features = $%d", argIndex))
		args = append(args, req.RateToolsFeatures)
		argIndex++
	}
	if req.EdocsAvailable != nil {
		setParts = append(setParts, fmt.Sprintf("edocs_available = $%d", argIndex))
		args = append(args, *req.EdocsAvailable)
		argIndex++
	}
	if req.EdocsLimit != nil {
		setParts = append(setParts, fmt.Sprintf("edocs_limit = $%d", argIndex))
		args = append(args, *req.EdocsLimit)
		argIndex++
	}
	if req.EdocsHasArchiving != nil {
		setParts = append(setParts, fmt.Sprintf("edocs_has_archiving = $%d", argIndex))
		args = append(args, *req.EdocsHasArchiving)
		argIndex++
	}
	if req.SupportLevel != nil {
		setParts = append(setParts, fmt.Sprintf("support_level = $%d", argIndex))
		args = append(args, *req.SupportLevel)
		argIndex++
	}
	if req.PaymentGuarantee != nil {
		setParts = append(setParts, fmt.Sprintf("payment_guarantee = $%d", argIndex))
		args = append(args, *req.PaymentGuarantee)
		argIndex++
	}
	if req.APIAccess != nil {
		setParts = append(setParts, fmt.Sprintf("api_access = $%d", argIndex))
		args = append(args, *req.APIAccess)
		argIndex++
	}
	if req.DisplayOrder != nil {
		setParts = append(setParts, fmt.Sprintf("display_order = $%d", argIndex))
		args = append(args, *req.DisplayOrder)
		argIndex++
	}
	if req.IsPopular != nil {
		setParts = append(setParts, fmt.Sprintf("is_popular = $%d", argIndex))
		args = append(args, *req.IsPopular)
		argIndex++
	}
	if req.IsRecommended != nil {
		setParts = append(setParts, fmt.Sprintf("is_recommended = $%d", argIndex))
		args = append(args, *req.IsRecommended)
		argIndex++
	}
	if req.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, *req.Description)
		argIndex++
	}
	if req.FeaturesSummary != nil {
		setParts = append(setParts, fmt.Sprintf("features_summary = $%d", argIndex))
		args = append(args, *req.FeaturesSummary)
		argIndex++
	}
	if req.AvailableFrom != nil {
		setParts = append(setParts, fmt.Sprintf("available_from = $%d", argIndex))
		args = append(args, *req.AvailableFrom)
		argIndex++
	}
	if req.AvailableUntil != nil {
		setParts = append(setParts, fmt.Sprintf("available_until = $%d", argIndex))
		args = append(args, *req.AvailableUntil)
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

	args = append(args, id)
	setClause := strings.Join(setParts, ", ")
	query := fmt.Sprintf(`
		UPDATE tbl_plan 
		SET %s 
		WHERE uuid = $%d AND deleted = 0 
		RETURNING id, uuid, name, code, provider, region, price_usd, price_local, 
		          local_currency, billing_cycle, load_posts_limit, load_posts_unlimited,
		          gps_tracking_level, gps_has_eta, rate_tools_level, rate_tools_features,
		          edocs_available, edocs_limit, edocs_has_archiving, support_level,
		          payment_guarantee, api_access, display_order, is_popular, is_recommended,
		          description, features_summary, available_from, available_until,
		          meta, meta2, meta3, created_at, updated_at`, setClause, argIndex)

	err := pgxscan.Get(context.Background(), database.DB, &plan, query, args...)
	return plan, err
}

func DeletePlanInternal(id uuid.UUID) error {
	query := `UPDATE tbl_plan SET deleted = 1, updated_at = NOW() WHERE uuid = $1 AND deleted = 0`
	_, err := database.DB.Exec(context.Background(), query, id)
	return err
}
