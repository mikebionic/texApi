package services

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"net/http"
	"texApi/database"
	"texApi/internal/dto"
	"texApi/pkg/utils"
)

// GetAnalytics handles GET request for analytics data
func GetAnalytics(ctx *gin.Context) {
	var filter dto.AnalyticsFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid filter parameters", err.Error()))
		return
	}

	// Set default sorting
	if filter.OrderBy == "" {
		filter.OrderBy = "created_at"
	}
	if filter.OrderDir == "" {
		filter.OrderDir = "desc"
	}

	analytics, total, err := GetAnalyticsList(filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to retrieve analytics", err.Error()))
		return
	}

	stats, err := GetAnalyticsStats(filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to retrieve analytics stats", err.Error()))
		return
	}

	response := struct {
		Total   int                `json:"total"`
		Page    int                `json:"page"`
		PerPage int                `json:"per_page"`
		Stats   dto.AnalyticsStats `json:"stats"`
		Data    []dto.Analytics    `json:"data"`
	}{
		Total:   total,
		Page:    filter.Page,
		PerPage: filter.PerPage,
		Stats:   stats,
		Data:    analytics,
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Analytics", response))
}

// GetAnalyticsList retrieves analytics data with filtering and pagination
func GetAnalyticsList(filter dto.AnalyticsFilter) ([]dto.Analytics, int, error) {
	var analytics []dto.Analytics
	var total int

	whereClauses := []string{"deleted = 0"}
	args := []interface{}{}
	argIndex := 1

	// Date range filters
	if filter.DateFrom != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("created_at >= $%d", argIndex))
		args = append(args, *filter.DateFrom)
		argIndex++
	}
	if filter.DateTo != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("created_at <= $%d", argIndex))
		args = append(args, *filter.DateTo)
		argIndex++
	}

	// Period filters
	if filter.PeriodStart != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("period_start >= $%d", argIndex))
		args = append(args, *filter.PeriodStart)
		argIndex++
	}
	if filter.PeriodEnd != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("period_end <= $%d", argIndex))
		args = append(args, *filter.PeriodEnd)
		argIndex++
	}

	// Value range filters
	if filter.UserAllMin != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("user_all >= $%d", argIndex))
		args = append(args, *filter.UserAllMin)
		argIndex++
	}
	if filter.UserAllMax != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("user_all <= $%d", argIndex))
		args = append(args, *filter.UserAllMax)
		argIndex++
	}

	// User type filters
	if filter.UserSenderMin != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("user_sender >= $%d", argIndex))
		args = append(args, *filter.UserSenderMin)
		argIndex++
	}
	if filter.UserSenderMax != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("user_sender <= $%d", argIndex))
		args = append(args, *filter.UserSenderMax)
		argIndex++
	}
	if filter.UserCarrierMin != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("user_carrier >= $%d", argIndex))
		args = append(args, *filter.UserCarrierMin)
		argIndex++
	}
	if filter.UserCarrierMax != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("user_carrier <= $%d", argIndex))
		args = append(args, *filter.UserCarrierMax)
		argIndex++
	}

	// Offer filters
	if filter.OfferAllMin != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("offer_all >= $%d", argIndex))
		args = append(args, *filter.OfferAllMin)
		argIndex++
	}
	if filter.OfferAllMax != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("offer_all <= $%d", argIndex))
		args = append(args, *filter.OfferAllMax)
		argIndex++
	}
	if filter.OfferActiveMin != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("offer_active >= $%d", argIndex))
		args = append(args, *filter.OfferActiveMin)
		argIndex++
	}
	if filter.OfferActiveMax != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("offer_active <= $%d", argIndex))
		args = append(args, *filter.OfferActiveMax)
		argIndex++
	}

	// Revenue filters
	if filter.RevenueMin != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("total_revenue >= $%d", argIndex))
		args = append(args, *filter.RevenueMin)
		argIndex++
	}
	if filter.RevenueMax != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("total_revenue <= $%d", argIndex))
		args = append(args, *filter.RevenueMax)
		argIndex++
	}

	whereClause := strings.Join(whereClauses, " AND ")

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM tbl_analytics WHERE %s", whereClause)
	err := database.DB.QueryRow(context.Background(), countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Main query with sorting and pagination
	query := fmt.Sprintf(`
        SELECT id, uuid, user_all, user_sender, user_carrier, last_user_id,
               user_sender_new, user_carrier_new, last_offer_id, offer_new_sender,
               offer_new_carrier, offer_all, offer_active, offer_pending,
               offer_completed, offer_no_response, last_completed_offer_id,
               total_revenue, average_cost_per_km, total_distance, active_companies,
               period_start, period_end, created_at, updated_at
        FROM tbl_analytics
        WHERE %s
        ORDER BY %s %s
        LIMIT $%d OFFSET $%d`,
		whereClause, filter.OrderBy, filter.OrderDir, argIndex, argIndex+1)

	args = append(args, filter.PerPage, (filter.Page-1)*filter.PerPage)
	err = pgxscan.Select(context.Background(), database.DB, &analytics, query, args...)
	return analytics, total, err
}

// GetAnalyticsStats retrieves summary statistics
func GetAnalyticsStats(filter dto.AnalyticsFilter) (dto.AnalyticsStats, error) {
	var stats dto.AnalyticsStats

	whereClauses := []string{"deleted = 0"}
	args := []interface{}{}
	argIndex := 1

	// Apply same filters as main query
	if filter.DateFrom != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("created_at >= $%d", argIndex))
		args = append(args, *filter.DateFrom)
		argIndex++
	}
	if filter.DateTo != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("created_at <= $%d", argIndex))
		args = append(args, *filter.DateTo)
		argIndex++
	}

	whereClause := strings.Join(whereClauses, " AND ")

	query := fmt.Sprintf(`
        SELECT 
            COUNT(*) as total_records,
            AVG(user_all) as avg_users_per_period,
            AVG(offer_all) as avg_offers_per_period,
            SUM(total_revenue) as total_revenue,
            MAX(created_at) as last_update
        FROM tbl_analytics 
        WHERE %s`, whereClause)

	err := database.DB.QueryRow(context.Background(), query, args...).Scan(
		&stats.TotalRecords,
		&stats.AvgUsersPerPeriod,
		&stats.AvgOffersPerPeriod,
		&stats.TotalRevenue,
		&stats.LastUpdate,
	)

	if err != nil {
		return stats, err
	}

	// Calculate growth rate (simplified - last vs previous period)
	growthQuery := fmt.Sprintf(`
        WITH ordered_analytics AS (
            SELECT user_all, ROW_NUMBER() OVER (ORDER BY created_at DESC) as rn
            FROM tbl_analytics 
            WHERE %s
            LIMIT 2
        )
        SELECT 
            CASE 
                WHEN COUNT(*) = 2 THEN 
                    ROUND(((MAX(CASE WHEN rn = 1 THEN user_all END) - MAX(CASE WHEN rn = 2 THEN user_all END)) * 100.0 / 
                           NULLIF(MAX(CASE WHEN rn = 2 THEN user_all END), 0)), 2)
                ELSE 0
            END as growth_rate
        FROM ordered_analytics`, whereClause)

	err = database.DB.QueryRow(context.Background(), growthQuery, args...).Scan(&stats.GrowthRate)
	return stats, err
}

// GenerateAnalytics creates new analytics entry (called by scheduler)
func GenerateAnalytics() error {
	log.Println("Starting analytics generation...")

	var analytics dto.Analytics
	var lastAnalytics dto.Analytics

	// Get the last analytics record to determine baseline
	lastQuery := `
        SELECT COALESCE(MAX(last_user_id), 0) as last_user_id,
               COALESCE(MAX(last_offer_id), 0) as last_offer_id,
               COALESCE(MAX(last_completed_offer_id), 0) as last_completed_offer_id
        FROM tbl_analytics 
        WHERE deleted = 0`

	err := database.DB.QueryRow(context.Background(), lastQuery).Scan(
		&lastAnalytics.LastUserID,
		&lastAnalytics.LastOfferID,
		&lastAnalytics.LastCompletedOfferID,
	)

	if err != nil {
		log.Printf("Error getting last analytics: %v", err)
		// If no previous analytics, start from 0
		lastAnalytics = dto.Analytics{}
	}

	now := time.Now()
	periodStart := now.Add(-24 * time.Hour) // Last 24 hours

	// Generate all metrics
	analytics.PeriodStart = periodStart
	analytics.PeriodEnd = now

	// User metrics
	analytics.UserAll = getUserCount("")
	analytics.UserSender = getUserCount("sender")
	analytics.UserCarrier = getUserCount("carrier")
	analytics.LastUserID = getLastUserID()
	analytics.UserSenderNew = getNewUserCount("sender", lastAnalytics.LastUserID)
	analytics.UserCarrierNew = getNewUserCount("carrier", lastAnalytics.LastUserID)

	// Offer metrics
	analytics.LastOfferID = getLastOfferID()
	analytics.OfferNewSender = getNewOfferCount("sender", lastAnalytics.LastOfferID)
	analytics.OfferNewCarrier = getNewOfferCount("carrier", lastAnalytics.LastOfferID)
	analytics.OfferAll = getOfferCount("active")
	analytics.OfferActive = getOfferCountByState("active", "enabled", "working")
	analytics.OfferPending = getOfferCountByState("pending")
	analytics.OfferCompleted = getOfferCountByState("completed", "archived")
	analytics.OfferNoResponse = getOffersWithoutResponse()
	analytics.LastCompletedOfferID = getLastCompletedOfferID()

	// Additional metrics
	analytics.TotalRevenue = getTotalRevenue()
	analytics.AverageCostPerKm = getAverageCostPerKm()
	analytics.TotalDistance = getTotalDistance()
	analytics.ActiveCompanies = getActiveCompanies()

	// Insert the analytics record
	insertQuery := `
        INSERT INTO tbl_analytics (
            user_all, user_sender, user_carrier, last_user_id,
            user_sender_new, user_carrier_new, last_offer_id,
            offer_new_sender, offer_new_carrier, offer_all,
            offer_active, offer_pending, offer_completed,
            offer_no_response, last_completed_offer_id,
            total_revenue, average_cost_per_km, total_distance,
            active_companies, period_start, period_end
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
            $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21
        )`

	_, err = database.DB.Exec(context.Background(), insertQuery,
		analytics.UserAll, analytics.UserSender, analytics.UserCarrier, analytics.LastUserID,
		analytics.UserSenderNew, analytics.UserCarrierNew, analytics.LastOfferID,
		analytics.OfferNewSender, analytics.OfferNewCarrier, analytics.OfferAll,
		analytics.OfferActive, analytics.OfferPending, analytics.OfferCompleted,
		analytics.OfferNoResponse, analytics.LastCompletedOfferID,
		analytics.TotalRevenue, analytics.AverageCostPerKm, analytics.TotalDistance,
		analytics.ActiveCompanies, analytics.PeriodStart, analytics.PeriodEnd,
	)

	if err != nil {
		log.Printf("Error inserting analytics: %v", err)
		return err
	}

	// Update last run time
	updateConfigQuery := `
        UPDATE tbl_analytics_config 
        SET value = $1, updated_at = CURRENT_TIMESTAMP 
        WHERE key = 'last_analytics_run'`

	_, err = database.DB.Exec(context.Background(), updateConfigQuery, now.Format(time.RFC3339))
	if err != nil {
		log.Printf("Error updating last run time: %v", err)
	}

	log.Println("Analytics generation completed successfully")
	return nil
}

// Helper functions for metrics calculation
func getUserCount(role string) int {
	var count int
	query := "SELECT COUNT(*) FROM tbl_user WHERE deleted = 0 AND active = 1"
	if role != "" {
		query += fmt.Sprintf(" AND role = '%s'", role)
	}
	database.DB.QueryRow(context.Background(), query).Scan(&count)
	return count
}

func getLastUserID() int {
	var id int
	query := "SELECT COALESCE(MAX(id), 0) FROM tbl_user WHERE deleted = 0"
	database.DB.QueryRow(context.Background(), query).Scan(&id)
	return id
}

func getNewUserCount(role string, lastID int) int {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM tbl_user WHERE deleted = 0 AND active = 1 AND role = '%s' AND id > %d", role, lastID)
	database.DB.QueryRow(context.Background(), query).Scan(&count)
	return count
}

func getLastOfferID() int {
	var id int
	query := "SELECT COALESCE(MAX(id), 0) FROM tbl_offer WHERE deleted = 0"
	database.DB.QueryRow(context.Background(), query).Scan(&id)
	return id
}

func getNewOfferCount(role string, lastID int) int {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM tbl_offer WHERE deleted = 0 AND offer_role = '%s' AND id > %d", role, lastID)
	database.DB.QueryRow(context.Background(), query).Scan(&count)
	return count
}

func getOfferCount(excludeState string) int {
	var count int
	query := "SELECT COUNT(*) FROM tbl_offer WHERE deleted = 0 AND offer_state NOT IN ('deleted', 'pending', 'disabled')"
	database.DB.QueryRow(context.Background(), query).Scan(&count)
	return count
}

func getOfferCountByState(states ...string) int {
	var count int
	stateStr := "'" + strings.Join(states, "','") + "'"
	query := fmt.Sprintf("SELECT COUNT(*) FROM tbl_offer WHERE deleted = 0 AND offer_state IN (%s)", stateStr)
	database.DB.QueryRow(context.Background(), query).Scan(&count)
	return count
}

func getOffersWithoutResponse() int {
	var count int
	query := `
        SELECT COUNT(*) 
        FROM tbl_offer o 
        LEFT JOIN tbl_offer_response r ON o.id = r.offer_id AND r.deleted = 0
        WHERE o.deleted = 0 AND r.id IS NULL`
	database.DB.QueryRow(context.Background(), query).Scan(&count)
	return count
}

func getLastCompletedOfferID() int {
	var id int
	query := "SELECT COALESCE(MAX(id), 0) FROM tbl_offer WHERE deleted = 0 AND offer_state IN ('completed', 'archived')"
	database.DB.QueryRow(context.Background(), query).Scan(&id)
	return id
}

func getTotalRevenue() float64 {
	var revenue float64
	query := "SELECT COALESCE(SUM(cost_per_km * distance), 0) FROM tbl_offer WHERE deleted = 0 AND offer_state = 'completed'"
	database.DB.QueryRow(context.Background(), query).Scan(&revenue)
	return revenue
}

func getAverageCostPerKm() float64 {
	var avg float64
	query := "SELECT COALESCE(AVG(cost_per_km), 0) FROM tbl_offer WHERE deleted = 0 AND cost_per_km > 0"
	database.DB.QueryRow(context.Background(), query).Scan(&avg)
	return avg
}

func getTotalDistance() int {
	var distance int
	query := "SELECT COALESCE(SUM(distance), 0) FROM tbl_offer WHERE deleted = 0 AND offer_state = 'completed'"
	database.DB.QueryRow(context.Background(), query).Scan(&distance)
	return distance
}

func getActiveCompanies() int {
	var count int
	query := "SELECT COUNT(DISTINCT company_id) FROM tbl_offer WHERE deleted = 0 AND offer_state IN ('active', 'working')"
	database.DB.QueryRow(context.Background(), query).Scan(&count)
	return count
}
