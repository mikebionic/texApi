package dto

import (
	"time"
)

type Claim struct {
	ID                  int       `json:"id"`
	UUID                string    `json:"uuid"`
	UserID              *int      `json:"user_id"`
	CompanyID           *int      `json:"company_id"`
	Name                *string   `json:"name"`
	Email               *string   `json:"email"`
	Phone               *string   `json:"phone"`
	Address             *string   `json:"address"`
	CompanyName         *string   `json:"company_name"`
	Subject             *string   `json:"subject"`
	Description         *string   `json:"description"`
	AdditionalDetails   *string   `json:"additional_details"`
	ResponseTitle       *string   `json:"response_title"`
	ResponseDescription *string   `json:"response_description"`
	ClaimType           string    `json:"claim_type"`
	ClaimStatus         string    `json:"claim_status"`
	UrgencyLevel        *int      `json:"urgency_level"`
	Meta                *string   `json:"meta"`
	Meta2               *string   `json:"meta2"`
	Meta3               *string   `json:"meta3"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	Active              int       `json:"active"`
	Deleted             int       `json:"deleted"`
}

type NewClaimRequest struct {
	UserID            *int    `json:"user_id"`
	CompanyID         *int    `json:"company_id"`
	Name              *string `json:"name"`
	Email             *string `json:"email"`
	Phone             *string `json:"phone"`
	Address           *string `json:"address"`
	CompanyName       *string `json:"company_name"`
	Subject           *string `json:"subject"`
	Description       *string `json:"description"`
	AdditionalDetails *string `json:"additional_details"`
	ClaimType         string  `json:"claim_type"`
	UrgencyLevel      *int    `json:"urgency_level"`
	Meta              *string `json:"meta"`
	Meta2             *string `json:"meta2"`
	Meta3             *string `json:"meta3"`
}

type UpdateClaimRequest struct {
	ID                  int     `json:"id"`
	UserID              *int    `json:"user_id"`
	CompanyID           *int    `json:"company_id"`
	Name                *string `json:"name"`
	Email               *string `json:"email"`
	Phone               *string `json:"phone"`
	Address             *string `json:"address"`
	CompanyName         *string `json:"company_name"`
	Subject             *string `json:"subject"`
	Description         *string `json:"description"`
	AdditionalDetails   *string `json:"additional_details"`
	ResponseTitle       *string `json:"response_title"`
	ResponseDescription *string `json:"response_description"`
	ClaimType           *string `json:"claim_type"`
	ClaimStatus         *string `json:"claim_status"`
	UrgencyLevel        *int    `json:"urgency_level"`
	Meta                *string `json:"meta"`
	Meta2               *string `json:"meta2"`
	Meta3               *string `json:"meta3"`
	Active              *int    `json:"active"`
}

type ClaimFilter struct {
	ClaimType    *string `form:"claim_type"`
	ClaimStatus  *string `form:"claim_status"`
	UrgencyLevel *int    `form:"urgency_level"`
	UserID       *int    `form:"user_id"`
	CompanyID    *int    `form:"company_id"`
	Email        *string `form:"email"`
	Active       *int    `form:"active"`
	Search       *string `form:"search"`
	Page         int     `form:"page,default=1"`
	PerPage      int     `form:"per_page,default=10"`
}

type DeleteClaimRequest struct {
	ID int `json:"id"`
}
