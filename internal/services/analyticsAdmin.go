package services

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"texApi/database"
	"texApi/internal/dto"
	"texApi/pkg/utils"
)

// GetAnalyticsStatus returns the current status of analytics system
func GetAnalyticsStatus(ctx *gin.Context) {
	status := make(map[string]interface{})

	// Get configuration
	config, err := getAnalyticsConfig()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to get analytics config", err.Error()))
		return
	}

	status["config"] = config

	// Get last analytics record
	var lastRecord dto.Analytics
	query := `
        SELECT id, created_at, period_start, period_end, user_all, offer_all, total_revenue
        FROM tbl_analytics 
        WHERE deleted = 0 
        ORDER BY created_at DESC 
        LIMIT 1`

	err = database.DB.QueryRow(context.Background(), query).Scan(
		&lastRecord.ID,
		&lastRecord.CreatedAt,
		&lastRecord.PeriodStart,
		&lastRecord.PeriodEnd,
		&lastRecord.UserAll,
		&lastRecord.OfferAll,
		&lastRecord.TotalRevenue,
	)

	if err == nil {
		status["last_record"] = lastRecord
	}

	// Get total records count
	var totalRecords int
	countQuery := "SELECT COUNT(*) FROM tbl_analytics WHERE deleted = 0"
	database.DB.QueryRow(context.Background(), countQuery).Scan(&totalRecords)
	status["total_records"] = totalRecords

	// Calculate next run time
	if lastRun, err := getLastRunTime(); err == nil {
		if intervalDays, err := getLogIntervalDays(); err == nil {
			nextRun := lastRun.Add(time.Duration(intervalDays) * 24 * time.Hour)
			status["next_run"] = nextRun
			status["should_run_now"] = time.Now().After(nextRun)
		}
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Analytics Status", status))
}

// ForceGenerateAnalytics manually triggers analytics generation
func ForceGenerateAnalytics(ctx *gin.Context) {
	// Check if user has admin privileges (implement based on your auth system)
	// userRole := ctx.GetString("user_role")
	// if userRole != "admin" {
	//     ctx.JSON(http.StatusForbidden, utils.FormatErrorResponse("Insufficient privileges", ""))
	//     return
	// }

	err := GenerateAnalytics()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to generate analytics", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Analytics generated successfully", nil))
}

// UpdateAnalyticsConfig updates analytics configuration
func UpdateAnalyticsConfig(ctx *gin.Context) {
	var req struct {
		LogIntervalDays *int  `json:"log_interval_days" binding:"omitempty,min=1,max=365"`
		Enabled         *bool `json:"enabled"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request data", err.Error()))
		return
	}

	// Update log interval if provided
	if req.LogIntervalDays != nil {
		err := updateConfigValue("log_interval_days", strconv.Itoa(*req.LogIntervalDays))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to update log interval", err.Error()))
			return
		}
	}

	// Update enabled status if provided
	if req.Enabled != nil {
		enabledStr := "false"
		if *req.Enabled {
			enabledStr = "true"
		}
		err := updateConfigValue("enabled", enabledStr)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to update enabled status", err.Error()))
			return
		}
	}

	// Get updated configuration
	config, err := getAnalyticsConfig()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to get updated config", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Configuration updated successfully", config))
}

// GetAnalyticsConfig returns current analytics configuration
func GetAnalyticsConfig(ctx *gin.Context) {
	config, err := getAnalyticsConfig()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to get analytics config", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Analytics Configuration", config))
}

// Helper functions
func getAnalyticsConfig() (map[string]interface{}, error) {
	config := make(map[string]interface{})

	query := "SELECT key, value, description FROM tbl_analytics_config"
	rows, err := database.DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var key, value, description string
		if err := rows.Scan(&key, &value, &description); err != nil {
			continue
		}

		config[key] = map[string]string{
			"value":       value,
			"description": description,
		}
	}

	return config, nil
}

func updateConfigValue(key, value string) error {
	query := `
        UPDATE tbl_analytics_config 
        SET value = $1, updated_at = CURRENT_TIMESTAMP 
        WHERE key = $2`

	_, err := database.DB.Exec(context.Background(), query, value, key)
	return err
}

func getLastRunTime() (time.Time, error) {
	var lastRunStr string
	query := "SELECT value FROM tbl_analytics_config WHERE key = 'last_analytics_run'"

	err := database.DB.QueryRow(context.Background(), query).Scan(&lastRunStr)
	if err != nil {
		return time.Time{}, err
	}

	return time.Parse(time.RFC3339, lastRunStr)
}

func getLogIntervalDays() (int, error) {
	var intervalDays string
	query := "SELECT value FROM tbl_analytics_config WHERE key = 'log_interval_days'"

	err := database.DB.QueryRow(context.Background(), query).Scan(&intervalDays)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(intervalDays)
}
