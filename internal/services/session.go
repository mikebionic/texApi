package services

import (
	"net/http"
	"strconv"
	"texApi/internal/dto"
	"texApi/internal/repo"
	"texApi/pkg/utils"

	"github.com/gin-gonic/gin"
)

func ExtractDeviceInfo(ctx *gin.Context) (string, string, string, string, string) {
	deviceName := ctx.GetHeader("X-Device-Name")
	deviceModel := ctx.GetHeader("X-Device-Model")
	deviceFirmware := ctx.GetHeader("X-Device-Firmware")
	appName := ctx.GetHeader("X-App-Name")
	appVersion := ctx.GetHeader("X-App-Version")

	return deviceName, deviceModel, deviceFirmware, appName, appVersion
}

func ListUserSessions(ctx *gin.Context) {
	userID := ctx.MustGet("id").(int)

	sessions, err := repo.GetUserSessions(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error retrieving sessions", "Database error"))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Active sessions", sessions))
}

func ListAllSessions(ctx *gin.Context) {
	role := ctx.MustGet("role").(string)
	if role != "admin" {
		ctx.JSON(http.StatusForbidden, utils.FormatErrorResponse("Permission denied", "Admin access required"))
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))
	orderBy := ctx.DefaultQuery("order_by", "created_at")
	orderDir := ctx.DefaultQuery("order_dir", "DESC")

	userID := ctx.Query("user_id")
	var userIDPtr *int
	if userID != "" {
		id, err := strconv.Atoi(userID)
		if err == nil {
			userIDPtr = &id
		}
	}

	companyID := ctx.Query("company_id")
	var companyIDPtr *int
	if companyID != "" {
		id, err := strconv.Atoi(companyID)
		if err == nil {
			companyIDPtr = &id
		}
	}

	loginMethod := ctx.Query("login_method")
	var loginMethodPtr *string
	if loginMethod != "" {
		loginMethodPtr = &loginMethod
	}

	isActiveStr := ctx.Query("is_active")
	var isActivePtr *bool
	if isActiveStr != "" {
		isActive := isActiveStr == "true"
		isActivePtr = &isActive
	}

	params := dto.SessionListParams{
		UserID:      userIDPtr,
		CompanyID:   companyIDPtr,
		LoginMethod: loginMethodPtr,
		IsActive:    isActivePtr,
		Page:        page,
		PerPage:     perPage,
		OrderBy:     orderBy,
		OrderDir:    orderDir,
	}

	sessions, totalCount, err := repo.ListSessions(params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error retrieving sessions", "Database error"))
		return
	}

	response := utils.PaginatedResponse{
		Total:   totalCount,
		Page:    page,
		PerPage: perPage,
		Data:    sessions,
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Sessions", response))
}

func RevokeSession(ctx *gin.Context) {
	userID := ctx.MustGet("id").(int)
	sessionID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid session ID", ""))
		return
	}

	session, err := repo.GetSessionByID(sessionID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Session not found", ""))
		return
	}

	if session.UserID != userID {
		ctx.JSON(http.StatusForbidden, utils.FormatErrorResponse("Access denied", ""))
		return
	}

	err = repo.InvalidateSessionByID(sessionID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error revoking session", ""))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Session revoked successfully", nil))
}

func AdminRevokeSession(ctx *gin.Context) {
	sessionID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid session ID", ""))
		return
	}

	err = repo.InvalidateSessionByID(sessionID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error revoking session", ""))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Session revoked successfully", nil))
}

func AdminRevokeUserSessions(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid user ID", ""))
		return
	}

	user, err := repo.GetUserById(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("User not found", ""))
		return
	}

	err = repo.InvalidateAllUserSessions(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Error revoking sessions", ""))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("All sessions revoked for user "+user.Username, nil))
}
