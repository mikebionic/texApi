package dto

import (
	"github.com/google/uuid"
)

type User struct {
	ID                     int       `json:"id"`
	UUID                   uuid.UUID `json:"uuid"`
	Username               string    `json:"username"`
	Password               string    `json:"-"`
	Email                  string    `json:"email"`
	Phone                  string    `json:"phone"`
	Role                   string    `json:"role"`
	RoleID                 int       `json:"role_id"`
	CompanyID              int       `json:"company_id"`
	Verified               int       `json:"verified"`
	CreatedAt              string    `json:"created_at"`
	UpdatedAt              string    `json:"updated_at"`
	Active                 int       `json:"active"`
	Deleted                int       `json:"deleted"`
	OauthProvider          string    `json:"oauth_provider"`
	OauthUserID            string    `json:"oauth_user_id"`
	OauthLocation          string    `json:"oauth_location"`
	OauthAccessToken       string    `json:"oauth_access_token"`
	OauthAccessTokenSecret string    `json:"oauth_access_token_secret"`
	OauthRefreshToken      string    `json:"oauth_refresh_token"`
	OauthExpiresAt         string    `json:"oauth_expires_at"`
	OauthIDToken           string    `json:"oauth_id_token"`
	RefreshToken           string    `json:"refresh_token"`
	VerifyTime             string    `json:"verify_time"`
	OTPKey                 string    `json:"otp_key"`
}

type CreateUser struct {
	Username               string `json:"username,omitempty"`
	Password               string `json:"password,omitempty"`
	Email                  string `json:"email,omitempty"`
	Phone                  string `json:"phone,omitempty"`
	Role                   string `json:"role,omitempty"`
	RoleID                 int    `json:"role_id,omitempty"`
	CompanyID              int    `json:"company_id,omitempty"`
	Verified               int    `json:"verified,omitempty"`
	Active                 int    `json:"active,omitempty"`
	OauthProvider          string `json:"oauth_provider,omitempty"`
	OauthUserID            string `json:"oauth_user_id,omitempty"`
	OauthLocation          string `json:"oauth_location,omitempty"`
	OauthAccessToken       string `json:"oauth_access_token,omitempty"`
	OauthAccessTokenSecret string `json:"oauth_access_token_secret,omitempty"`
	OauthRefreshToken      string `json:"oauth_refresh_token,omitempty"`
	OauthExpiresAt         string `json:"oauth_expires_at,omitempty"`
	OauthIDToken           string `json:"oauth_id_token,omitempty"`
	RefreshToken           string `json:"refresh_token,omitempty"`
	VerifyTime             string `json:"verify_time"`
}

type ProfileUpdate struct {
	Username *string `json:"username,omitempty"`
	Password *string `json:"password,omitempty"`
	Email    *string `json:"email,omitempty"`
	Phone    *string `json:"phone,omitempty"`
}

type RefreshTokenForm struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
