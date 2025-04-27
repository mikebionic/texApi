package dto

import (
	"time"
)

type UserDetails struct {
	ID        int    `json:"id" db:"id"`
	UUID      string `json:"uuid" db:"uuid"`
	Username  string `json:"username" db:"username"`
	Password  string `json:"-" db:"password"`
	Email     string `json:"email" db:"email"`
	Phone     string `json:"phone" db:"phone"`
	Role      string `json:"role" db:"role"`
	RoleID    int    `json:"role_id,omitempty"`
	CompanyID int    `json:"company_id" db:"company_id"`
	Verified  int    `json:"verified" db:"verified"`

	Meta  string `json:"meta"`
	Meta2 string `json:"meta2"`
	Meta3 string `json:"meta3"`

	RefreshToken string    `json:"refresh_token"`
	OTPKey       string    `json:"otp_key"`
	VerifyTime   time.Time `json:"verify_time"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	Active       int       `json:"active" db:"active"`
	Deleted      int       `json:"deleted" db:"deleted"`
	TotalCount   int       `json:"total_count" db:"total_count"`
}

type UserCreate struct {
	Username  string  `json:"username"`
	Email     string  `json:"email"`
	Phone     string  `json:"phone"`
	Password  string  `json:"password"`
	CompanyID *int    `json:"company_id,omitempty"`
	Role      string  `json:"role,omitempty"`
	RoleID    int     `json:"role_id,omitempty"`
	Meta      *string `json:"meta,omitempty"`
	Meta2     *string `json:"meta2,omitempty"`
	Meta3     *string `json:"meta3,omitempty"`
	Verified  *int    `json:"verified,omitempty"`

	// for company:
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
}

type UserUpdate struct {
	Username  *string `json:"username,omitempty"`
	Email     *string `json:"email,omitempty"`
	Phone     *string `json:"phone,omitempty"`
	Role      *string `json:"role,omitempty"`
	RoleID    *int    `json:"role_id,omitempty"`
	Active    *int    `json:"active,omitempty"`
	Verified  *int    `json:"verified,omitempty"`
	Password  *string `json:"password,omitempty"`
	Meta      *string `json:"meta,omitempty"`
	Meta2     *string `json:"meta2,omitempty"`
	Meta3     *string `json:"meta3,omitempty"`
	CompanyID *int    `json:"company_id,omitempty"`
}

type UserRichInfo struct {
	User    UserDetails    `json:"user"`
	Company CompanyDetails `json:"company"`
}
