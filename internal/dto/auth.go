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
	FirstName              string    `json:"first_name"`
	LastName               string    `json:"last_name"`
	NickName               string    `json:"nick_name"`
	AvatarURL              string    `json:"avatar_url"`
	Phone                  string    `json:"phone"`
	InfoPhone              string    `json:"info_phone"`
	Address                string    `json:"address"`
	RoleID                 int       `json:"role_id"`
	SubroleID              int       `json:"subrole_id"`
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
	OTPKey                 string    `json:"otp_key"`
	CheckTime              string    `json:"check_time"`
}

type CreateUser struct {
	Username               *string `json:"username,omitempty"`
	Password               *string `json:"password,omitempty"`
	Email                  *string `json:"email,omitempty"`
	FirstName              *string `json:"first_name,omitempty"`
	LastName               *string `json:"last_name,omitempty"`
	NickName               *string `json:"nick_name,omitempty"`
	AvatarURL              *string `json:"avatar_url,omitempty"`
	Phone                  *string `json:"phone,omitempty"`
	InfoPhone              *string `json:"info_phone,omitempty"`
	Address                *string `json:"address,omitempty"`
	RoleID                 int     `json:"role_id,omitempty" binding:"gt=0"`
	SubroleID              *int    `json:"subrole_id,omitempty"`
	Verified               int     `json:"verified,omitempty"`
	Active                 int     `json:"active,omitempty"`
	OauthProvider          *string `json:"oauth_provider,omitempty"`
	OauthUserID            *string `json:"oauth_user_id,omitempty"`
	OauthLocation          *string `json:"oauth_location,omitempty"`
	OauthAccessToken       *string `json:"oauth_access_token,omitempty"`
	OauthAccessTokenSecret *string `json:"oauth_access_token_secret,omitempty"`
	OauthRefreshToken      *string `json:"oauth_refresh_token,omitempty"`
	OauthIDToken           *string `json:"oauth_id_token,omitempty"`
}

type RefreshTokenForm struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
