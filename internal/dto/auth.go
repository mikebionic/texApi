package dto

import (
	"github.com/google/uuid"
)

type User struct {
	ID           int       `json:"id"`
	UUID         uuid.UUID `json:"uuid"`
	Username     string    `json:"username"`
	Password     string    `json:"-"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	Role         string    `json:"role"`
	RoleID       int       `json:"role_id"`
	CompanyID    int       `json:"company_id"`
	Verified     int       `json:"verified"`
	RefreshToken string    `json:"refresh_token"`
	OTPKey       string    `json:"otp_key"`
	VerifyTime   string    `json:"verify_time"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
	Active       int       `json:"active"`
	Deleted      int       `json:"deleted"`
}

type CreateUser struct {
	Username     string  `json:"username,omitempty"`
	Password     string  `json:"password,omitempty"`
	Email        string  `json:"email,omitempty"`
	Phone        string  `json:"phone,omitempty"`
	Role         string  `json:"role,omitempty"`
	RoleID       int     `json:"role_id,omitempty"`
	CompanyID    int     `json:"company_id,omitempty"`
	Verified     int     `json:"verified,omitempty"`
	Meta         string  `json:"meta"`
	Meta2        string  `json:"meta2"`
	Meta3        string  `json:"meta3"`
	OTP          *string `json:"otp_key"`
	RefreshToken string  `json:"refresh_token,omitempty"`
	VerifyTime   string  `json:"verify_time,omitempty"`
	Active       int     `json:"active,omitempty"`
}

type UserUpdateAuth struct {
	Username *string `json:"username,omitempty"`
	Password *string `json:"password,omitempty"`
	Email    *string `json:"email,omitempty"`
	Phone    *string `json:"phone,omitempty"`
}

type RefreshTokenForm struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type OAuthUser struct {
	RawData   map[string]interface{}
	Provider  string
	Email     string
	Name      string
	FirstName string
	LastName  string
	AvatarURL string
}
