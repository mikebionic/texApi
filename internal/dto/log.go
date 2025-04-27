package dto

import (
	"time"
)

// VerifyRequestDetails represents the data structure for a verification request
type VerifyRequestDetails struct {
	ID          int       `json:"id" db:"id"`
	UserID      int       `json:"user_id" db:"user_id"`
	CompanyID   int       `json:"company_id" db:"company_id"`
	Status      string    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Deleted     int       `json:"deleted" db:"deleted"`
	Username    *string   `json:"username" db:"username"`
	FirstName   *string   `json:"first_name" db:"first_name"`
	LastName    *string   `json:"last_name" db:"last_name"`
	CompanyName *string   `json:"company_name" db:"company_name"`
	Role        *string   `json:"role" db:"role"`
	TotalCount  int       `json:"total_count" db:"total_count"`
}

// VerifyRequestUpdate represents the data structure for updating a verification request
type VerifyRequestUpdate struct {
	Status *string `json:"status,omitempty"`
}

// PlanMoveDetails represents the data structure for a plan movement
type PlanMoveDetails struct {
	ID          int       `json:"id" db:"id"`
	UserID      int       `json:"user_id" db:"user_id"`
	CompanyID   int       `json:"company_id" db:"company_id"`
	Status      *string   `json:"status" db:"status"`
	ValidUntil  time.Time `json:"valid_until" db:"valid_until"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Deleted     int       `json:"deleted" db:"deleted"`
	Username    *string   `json:"username" db:"username"`
	FirstName   *string   `json:"first_name" db:"first_name"`
	LastName    *string   `json:"last_name" db:"last_name"`
	CompanyName *string   `json:"company_name" db:"company_name"`
	Role        *string   `json:"role" db:"role"`
	TotalCount  int       `json:"total_count" db:"total_count"`
}

// PlanMoveCreate represents the data structure for creating a plan movement
type PlanMoveCreate struct {
	UserID    int    `json:"user_id"`
	CompanyID int    `json:"company_id"`
	Status    string `json:"status,omitempty"`
}

// PlanMoveUpdate represents the data structure for updating a plan movement
type PlanMoveUpdate struct {
	Status     *string    `json:"status,omitempty"`
	ValidUntil *time.Time `json:"valid_until,omitempty"`
}

type UserLogDetails struct {
	ID          int       `json:"id" db:"id"`
	UserID      int       `json:"user_id" db:"user_id"`
	CompanyID   int       `json:"company_id" db:"company_id"`
	Username    *string   `json:"username" db:"username"`
	FirstName   *string   `json:"first_name" db:"first_name"`
	LastName    *string   `json:"last_name" db:"last_name"`
	CompanyName *string   `json:"company_name" db:"company_name"`
	Role        string    `json:"role" db:"role"`
	Action      string    `json:"action" db:"action"`
	Details     string    `json:"details" db:"details"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Deleted     int       `json:"deleted" db:"deleted"`
	TotalCount  int       `json:"total_count" db:"total_count"`
}

// UserLogCreate represents the data structure for creating a user log entry
type UserLogCreate struct {
	UserID    int    `json:"user_id"`
	CompanyID int    `json:"company_id"`
	Role      string `json:"role"`
	Action    string `json:"action"`
	Details   string `json:"details"`
}
