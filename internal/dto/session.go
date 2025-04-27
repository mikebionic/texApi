package dto

import (
	"time"
)

type Session struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	CompanyID      int       `json:"company_id"`
	RefreshToken   string    `json:"-"`
	ExpiresAt      time.Time `json:"expires_at"`
	DeviceName     string    `json:"device_name"`
	DeviceModel    string    `json:"device_model"`
	DeviceFirmware string    `json:"device_firmware"`
	AppName        string    `json:"app_name"`
	AppVersion     string    `json:"app_version"`
	UserAgent      string    `json:"user_agent"`
	IPAddress      string    `json:"ip_address"`
	LoginMethod    string    `json:"login_method"`
	LastUsedAt     time.Time `json:"last_used_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	IsActive       bool      `json:"is_active"`
}

type SessionListItem struct {
	Session
	Username   string `json:"username"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Role       string `json:"role"`
	TotalCount int    `json:"total_count"`
}

type CreateSessionInput struct {
	UserID         int       `json:"user_id"`
	CompanyID      int       `json:"company_id"`
	RefreshToken   string    `json:"refresh_token"`
	ExpiresAt      time.Time `json:"expires_at"`
	DeviceName     string    `json:"device_name"`
	DeviceModel    string    `json:"device_model"`
	DeviceFirmware string    `json:"device_firmware"`
	AppName        string    `json:"app_name"`
	AppVersion     string    `json:"app_version"`
	UserAgent      string    `json:"user_agent"`
	IPAddress      string    `json:"ip_address"`
	LoginMethod    string    `json:"login_method"`
}

type SessionListParams struct {
	UserID      *int       `json:"user_id"`
	CompanyID   *int       `json:"company_id"`
	LoginMethod *string    `json:"login_method"`
	DeviceName  *string    `json:"device_name"`
	AppName     *string    `json:"app_name"`
	IsActive    *bool      `json:"is_active"`
	CreatedFrom *time.Time `json:"created_from"`
	CreatedTo   *time.Time `json:"created_to"`
	Page        int        `json:"page"`
	PerPage     int        `json:"per_page"`
	OrderBy     string     `json:"order_by"`
	OrderDir    string     `json:"order_dir"`
}
