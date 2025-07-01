package dto

import (
	"github.com/google/uuid"
	"time"
)

type Version struct {
	ID                     int        `json:"id"`
	UUID                   uuid.UUID  `json:"uuid"`
	VersionNumber          string     `json:"version_number"`
	VersionCode            int        `json:"version_code"`
	Title                  string     `json:"title"`
	Description            *string    `json:"description"`
	Platform               string     `json:"platform"`
	MinimalPlatformVersion *string    `json:"minimal_platform_version"`
	DownloadURL            *string    `json:"download_url"`
	FileSize               *int64     `json:"file_size"`
	Checksum               *string    `json:"checksum"`
	Changelog              *string    `json:"changelog"`
	ReleaseNotes           *string    `json:"release_notes"`
	IsCriticalUpdate       bool       `json:"is_critical_update"`
	IsBeta                 bool       `json:"is_beta"`
	AutoUpdateEnabled      bool       `json:"auto_update_enabled"`
	RolloutPercentage      int        `json:"rollout_percentage"`
	ActiveAt               *time.Time `json:"active_at"`
	DeprecatedAt           *time.Time `json:"deprecated_at"`
	EndOfLifeAt            *time.Time `json:"end_of_life_at"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
}

type CreateVersionRequest struct {
	VersionNumber          string     `json:"version_number" binding:"required,max=50"`
	VersionCode            int        `json:"version_code" binding:"required,min=1"`
	Title                  string     `json:"title" binding:"required,max=200"`
	Description            *string    `json:"description" binding:"omitempty,max=1000"`
	Platform               string     `json:"platform" binding:"required,oneof=ios android web desktop"`
	MinimalPlatformVersion *string    `json:"minimal_platform_version" binding:"omitempty,max=50"`
	DownloadURL            *string    `json:"download_url" binding:"omitempty,max=500,url"`
	FileSize               *int64     `json:"file_size" binding:"omitempty,min=0"`
	Checksum               *string    `json:"checksum" binding:"omitempty,max=128"`
	Changelog              *string    `json:"changelog" binding:"omitempty"`
	ReleaseNotes           *string    `json:"release_notes" binding:"omitempty"`
	IsCriticalUpdate       bool       `json:"is_critical_update"`
	IsBeta                 bool       `json:"is_beta"`
	AutoUpdateEnabled      bool       `json:"auto_update_enabled"`
	RolloutPercentage      int        `json:"rollout_percentage" binding:"min=0,max=100"`
	ActiveAt               *time.Time `json:"active_at" binding:"omitempty"`
}

type UpdateVersionRequest struct {
	VersionNumber          *string    `json:"version_number" binding:"omitempty,max=50"`
	VersionCode            *int       `json:"version_code" binding:"omitempty,min=1"`
	Title                  *string    `json:"title" binding:"omitempty,max=200"`
	Description            *string    `json:"description" binding:"omitempty,max=1000"`
	Platform               *string    `json:"platform" binding:"omitempty,oneof=ios android web desktop"`
	MinimalPlatformVersion *string    `json:"minimal_platform_version" binding:"omitempty,max=50"`
	DownloadURL            *string    `json:"download_url" binding:"omitempty,max=500,url"`
	FileSize               *int64     `json:"file_size" binding:"omitempty,min=0"`
	Checksum               *string    `json:"checksum" binding:"omitempty,max=128"`
	Changelog              *string    `json:"changelog" binding:"omitempty"`
	ReleaseNotes           *string    `json:"release_notes" binding:"omitempty"`
	IsCriticalUpdate       *bool      `json:"is_critical_update" binding:"omitempty"`
	IsBeta                 *bool      `json:"is_beta" binding:"omitempty"`
	AutoUpdateEnabled      *bool      `json:"auto_update_enabled" binding:"omitempty"`
	RolloutPercentage      *int       `json:"rollout_percentage" binding:"omitempty,min=0,max=100"`
	ActiveAt               *time.Time `json:"active_at" binding:"omitempty"`
	DeprecatedAt           *time.Time `json:"deprecated_at" binding:"omitempty"`
	EndOfLifeAt            *time.Time `json:"end_of_life_at" binding:"omitempty"`
}

type VersionFilter struct {
	Platform *string `form:"platform" binding:"omitempty,oneof=ios android web desktop"`
	IsBeta   *bool   `form:"is_beta"`
	Active   *bool   `form:"active"`
	Search   *string `form:"search"`
	Page     int     `form:"page,default=1" binding:"min=1"`
	PerPage  int     `form:"per_page,default=10" binding:"min=1,max=100"`
}

type UpdateCheckResponse struct {
	HasUpdate     bool     `json:"has_update"`
	LatestVersion *Version `json:"latest_version,omitempty"`
	IsCritical    bool     `json:"is_critical"`
	CanAutoUpdate bool     `json:"can_auto_update"`
	ShouldUpdate  bool     `json:"should_update"`
	UpdateMessage string   `json:"update_message"`
	IsDeprecated  bool     `json:"is_deprecated"`
	IsEndOfLife   bool     `json:"is_end_of_life"`
}
