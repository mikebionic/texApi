package services

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"strings"
	"texApi/database"
	"texApi/internal/dto"
	"texApi/pkg/utils"
	"time"
)

func GetVersions(ctx *gin.Context) {
	var filter dto.VersionFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid filter parameters", err.Error()))
		return
	}

	versions, total, err := GetVersionsList(filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to retrieve versions", err.Error()))
		return
	}

	response := utils.PaginatedResponse{
		Total:   total,
		Page:    filter.Page,
		PerPage: filter.PerPage,
		Data:    versions,
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Versions", response))
}

func GetVersionByID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("uuid"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid version ID", err.Error()))
		return
	}

	version, err := GetVersionByIDInternal(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("Version not found", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Version retrieved", version))
}

func GetLatestVersion(ctx *gin.Context) {
	platform := ctx.Param("platform")
	if platform == "" {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Platform is required", ""))
		return
	}

	version, err := GetLatestVersionInternal(platform)
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.FormatErrorResponse("No version found for platform", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Latest version", version))
}

func CheckForUpdates(ctx *gin.Context) {
	platform := ctx.Param("platform")
	currentVersionStr := ctx.Param("current_version")

	if platform == "" || currentVersionStr == "" {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Platform and current version are required", ""))
		return
	}

	currentVersionCode, err := strconv.Atoi(currentVersionStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid current version format", err.Error()))
		return
	}

	updateCheck, err := CheckForUpdatesInternal(platform, currentVersionCode)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to check for updates", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Update check", updateCheck))
}

func CreateVersion(ctx *gin.Context) {
	var req dto.CreateVersionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request data", err.Error()))
		return
	}

	version, err := CreateVersionInternal(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to create version", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.FormatResponse("Version created successfully", version))
}

func UpdateVersion(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("uuid"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid version ID", err.Error()))
		return
	}

	var req dto.UpdateVersionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid request data", err.Error()))
		return
	}

	version, err := UpdateVersionInternal(id, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to update version", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Version updated successfully", version))
}

func DeleteVersion(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("uuid"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid version ID", err.Error()))
		return
	}

	if err := DeleteVersionInternal(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to delete version", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Version deleted successfully", nil))
}

func ActivateVersion(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("uuid"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid version ID", err.Error()))
		return
	}

	if err := ActivateVersionInternal(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to activate version", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Version activated successfully", nil))
}

func DeprecateVersion(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("uuid"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.FormatErrorResponse("Invalid version ID", err.Error()))
		return
	}

	if err := DeprecateVersionInternal(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.FormatErrorResponse("Failed to deprecate version", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.FormatResponse("Version deprecated successfully", nil))
}

func GetVersionsList(filter dto.VersionFilter) ([]dto.Version, int, error) {
	var versions []dto.Version
	var total int

	whereClauses := []string{"deleted = 0"}
	args := []interface{}{}
	argIndex := 1

	if filter.Platform != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("platform = $%d", argIndex))
		args = append(args, *filter.Platform)
		argIndex++
	}
	if filter.IsBeta != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("is_beta = $%d", argIndex))
		args = append(args, *filter.IsBeta)
		argIndex++
	}
	if filter.Active != nil {
		if *filter.Active {
			whereClauses = append(whereClauses, fmt.Sprintf("active = 1"))
		} else {
			whereClauses = append(whereClauses, fmt.Sprintf("active = 0"))
		}
	}
	if filter.Search != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("(title ILIKE $%d OR version_number ILIKE $%d OR description ILIKE $%d)", argIndex, argIndex, argIndex))
		args = append(args, "%"+*filter.Search+"%")
		argIndex++
	}

	whereClause := strings.Join(whereClauses, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM tbl_version WHERE %s", whereClause)
	err := database.DB.QueryRow(context.Background(), countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT id, uuid, version_number, version_code, title, description, platform,
		       minimal_platform_version, download_url, file_size, checksum, changelog,
		       release_notes, is_critical_update, is_beta, auto_update_enabled,
		       rollout_percentage, active_at, deprecated_at, end_of_life_at,
		       created_at, updated_at
		FROM tbl_version
		WHERE %s
		ORDER BY version_code DESC, created_at DESC
		LIMIT $%d OFFSET $%d`,
		whereClause, argIndex, argIndex+1)

	args = append(args, filter.PerPage, (filter.Page-1)*filter.PerPage)
	err = pgxscan.Select(context.Background(), database.DB, &versions, query, args...)
	return versions, total, err
}

func GetVersionByIDInternal(id uuid.UUID) (dto.Version, error) {
	var version dto.Version
	query := `
		SELECT id, uuid, version_number, version_code, title, description, platform,
		       minimal_platform_version, download_url, file_size, checksum, changelog,
		       release_notes, is_critical_update, is_beta, auto_update_enabled,
		       rollout_percentage, active_at, deprecated_at, end_of_life_at,
		       created_at, updated_at
		FROM tbl_version
		WHERE uuid = $1 AND deleted = 0`
	err := pgxscan.Get(context.Background(), database.DB, &version, query, id)
	return version, err
}

func GetLatestVersionInternal(platform string) (dto.Version, error) {
	var version dto.Version
	query := `
		SELECT id, uuid, version_number, version_code, title, description, platform,
		       minimal_platform_version, download_url, file_size, checksum, changelog,
		       release_notes, is_critical_update, is_beta, auto_update_enabled,
		       rollout_percentage, active_at, deprecated_at, end_of_life_at,
		       created_at, updated_at
		FROM tbl_version
		WHERE platform = $1 AND active = 1 AND deleted = 0 AND is_beta = false
		ORDER BY version_code DESC
		LIMIT 1`
	err := pgxscan.Get(context.Background(), database.DB, &version, query, platform)
	return version, err
}

func CheckForUpdatesInternal(platform string, currentVersionCode int) (dto.UpdateCheckResponse, error) {
	latestVersion, err := GetLatestVersionInternal(platform)
	if err != nil {
		return dto.UpdateCheckResponse{}, err
	}

	var currentVersion dto.Version
	currentQuery := `
		SELECT deprecated_at, end_of_life_at, is_critical_update
		FROM tbl_version
		WHERE platform = $1 AND version_code = $2 AND deleted = 0
		LIMIT 1`
	_ = pgxscan.Get(context.Background(), database.DB, &currentVersion, currentQuery, platform, currentVersionCode)

	now := time.Now()
	hasUpdate := latestVersion.VersionCode > currentVersionCode
	isCritical := hasUpdate && latestVersion.IsCriticalUpdate
	canAutoUpdate := hasUpdate && latestVersion.AutoUpdateEnabled
	isDeprecated := currentVersion.DeprecatedAt != nil && currentVersion.DeprecatedAt.Before(now)
	isEndOfLife := currentVersion.EndOfLifeAt != nil && currentVersion.EndOfLifeAt.Before(now)
	shouldUpdate := hasUpdate && (isCritical || isDeprecated || isEndOfLife)

	updateMessage := "Your app is up to date"
	if hasUpdate {
		if isCritical {
			updateMessage = "Critical update required. Please update immediately."
		} else if isEndOfLife {
			updateMessage = "Your version is no longer supported. Please update."
		} else if isDeprecated {
			updateMessage = "Your version is deprecated. Please consider updating."
		} else {
			updateMessage = "A new version is available."
		}
	}

	response := dto.UpdateCheckResponse{
		HasUpdate:     hasUpdate,
		IsCritical:    isCritical,
		CanAutoUpdate: canAutoUpdate,
		ShouldUpdate:  shouldUpdate,
		UpdateMessage: updateMessage,
		IsDeprecated:  isDeprecated,
		IsEndOfLife:   isEndOfLife,
	}

	if hasUpdate {
		response.LatestVersion = &latestVersion
	}

	return response, nil
}

func CreateVersionInternal(req dto.CreateVersionRequest) (dto.Version, error) {
	var version dto.Version
	query := `
		INSERT INTO tbl_version (
			version_number, version_code, title, description, platform,
			minimal_platform_version, download_url, file_size, checksum,
			changelog, release_notes, is_critical_update, is_beta,
			auto_update_enabled, rollout_percentage, active_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		) RETURNING id, uuid, version_number, version_code, title, description, platform,
		           minimal_platform_version, download_url, file_size, checksum, changelog,
		           release_notes, is_critical_update, is_beta, auto_update_enabled,
		           rollout_percentage, active_at, deprecated_at, end_of_life_at,
		           created_at, updated_at`

	err := pgxscan.Get(context.Background(), database.DB, &version, query,
		req.VersionNumber, req.VersionCode, req.Title, req.Description, req.Platform,
		req.MinimalPlatformVersion, req.DownloadURL, req.FileSize, req.Checksum,
		req.Changelog, req.ReleaseNotes, req.IsCriticalUpdate, req.IsBeta,
		req.AutoUpdateEnabled, req.RolloutPercentage, req.ActiveAt)
	return version, err
}

func UpdateVersionInternal(id uuid.UUID, req dto.UpdateVersionRequest) (dto.Version, error) {
	var version dto.Version
	setParts := []string{"updated_at = NOW()"}
	args := []interface{}{}
	argIndex := 1

	if req.VersionNumber != nil {
		setParts = append(setParts, fmt.Sprintf("version_number = $%d", argIndex))
		args = append(args, *req.VersionNumber)
		argIndex++
	}
	if req.VersionCode != nil {
		setParts = append(setParts, fmt.Sprintf("version_code = $%d", argIndex))
		args = append(args, *req.VersionCode)
		argIndex++
	}
	if req.Title != nil {
		setParts = append(setParts, fmt.Sprintf("title = $%d", argIndex))
		args = append(args, *req.Title)
		argIndex++
	}
	if req.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, *req.Description)
		argIndex++
	}
	if req.Platform != nil {
		setParts = append(setParts, fmt.Sprintf("platform = $%d", argIndex))
		args = append(args, *req.Platform)
		argIndex++
	}
	if req.MinimalPlatformVersion != nil {
		setParts = append(setParts, fmt.Sprintf("minimal_platform_version = $%d", argIndex))
		args = append(args, *req.MinimalPlatformVersion)
		argIndex++
	}
	if req.DownloadURL != nil {
		setParts = append(setParts, fmt.Sprintf("download_url = $%d", argIndex))
		args = append(args, *req.DownloadURL)
		argIndex++
	}
	if req.FileSize != nil {
		setParts = append(setParts, fmt.Sprintf("file_size = $%d", argIndex))
		args = append(args, *req.FileSize)
		argIndex++
	}
	if req.Checksum != nil {
		setParts = append(setParts, fmt.Sprintf("checksum = $%d", argIndex))
		args = append(args, *req.Checksum)
		argIndex++
	}
	if req.Changelog != nil {
		setParts = append(setParts, fmt.Sprintf("changelog = $%d", argIndex))
		args = append(args, *req.Changelog)
		argIndex++
	}
	if req.ReleaseNotes != nil {
		setParts = append(setParts, fmt.Sprintf("release_notes = $%d", argIndex))
		args = append(args, *req.ReleaseNotes)
		argIndex++
	}
	if req.IsCriticalUpdate != nil {
		setParts = append(setParts, fmt.Sprintf("is_critical_update = $%d", argIndex))
		args = append(args, *req.IsCriticalUpdate)
		argIndex++
	}
	if req.IsBeta != nil {
		setParts = append(setParts, fmt.Sprintf("is_beta = $%d", argIndex))
		args = append(args, *req.IsBeta)
		argIndex++
	}
	if req.AutoUpdateEnabled != nil {
		setParts = append(setParts, fmt.Sprintf("auto_update_enabled = $%d", argIndex))
		args = append(args, *req.AutoUpdateEnabled)
		argIndex++
	}
	if req.RolloutPercentage != nil {
		setParts = append(setParts, fmt.Sprintf("rollout_percentage = $%d", argIndex))
		args = append(args, *req.RolloutPercentage)
		argIndex++
	}
	if req.ActiveAt != nil {
		setParts = append(setParts, fmt.Sprintf("active_at = $%d", argIndex))
		args = append(args, *req.ActiveAt)
		argIndex++
	}
	if req.DeprecatedAt != nil {
		setParts = append(setParts, fmt.Sprintf("deprecated_at = $%d", argIndex))
		args = append(args, *req.DeprecatedAt)
		argIndex++
	}
	if req.EndOfLifeAt != nil {
		setParts = append(setParts, fmt.Sprintf("end_of_life_at = $%d", argIndex))
		args = append(args, *req.EndOfLifeAt)
		argIndex++
	}

	args = append(args, id)
	setClause := strings.Join(setParts, ", ")
	query := fmt.Sprintf(`
		UPDATE tbl_version 
		SET %s 
		WHERE uuid = $%d AND deleted = 0 
		RETURNING id, uuid, version_number, version_code, title, description, platform,
		          minimal_platform_version, download_url, file_size, checksum, changelog,
		          release_notes, is_critical_update, is_beta, auto_update_enabled,
		          rollout_percentage, active_at, deprecated_at, end_of_life_at,
		          created_at, updated_at`, setClause, argIndex)

	err := pgxscan.Get(context.Background(), database.DB, &version, query, args...)
	return version, err
}

func DeleteVersionInternal(id uuid.UUID) error {
	query := `UPDATE tbl_version SET deleted = 1, updated_at = NOW() WHERE uuid = $1 AND deleted = 0`
	_, err := database.DB.Exec(context.Background(), query, id)
	return err
}

func ActivateVersionInternal(id uuid.UUID) error {
	query := `UPDATE tbl_version SET active = 1, active_at = NOW(), updated_at = NOW() WHERE uuid = $1 AND deleted = 0`
	_, err := database.DB.Exec(context.Background(), query, id)
	return err
}

func DeprecateVersionInternal(id uuid.UUID) error {
	query := `UPDATE tbl_version SET deprecated_at = NOW(), updated_at = NOW() WHERE uuid = $1 AND deleted = 0`
	_, err := database.DB.Exec(context.Background(), query, id)
	return err
}
