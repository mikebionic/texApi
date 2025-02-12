package dto

import (
	"time"
)

type OfferResponseDetails struct {
	OfferResponse
	Company   *CompanyCreate `json:"company,omitempty" db:"company,omitempty"`
	ToCompany *CompanyCreate `json:"to_company,omitempty" db:"to_company,omitempty"`
	Offer     *Offer         `json:"offer,omitempty" db:"offer,omitempty"`
}

type OfferResponse struct {
	ID          int       `json:"id,omitempty" db:"id"`
	UUID        string    `json:"uuid,omitempty" db:"uuid"`
	CompanyID   int       `json:"company_id" validate:"required"`
	OfferID     int       `json:"offer_id" validate:"required"`
	ToCompanyID int       `json:"to_company_id" validate:"required"`
	State       string    `json:"state" validate:"required"`
	BidPrice    *float64  `json:"bid_price,omitempty"`
	Title       *string   `json:"title,omitempty"`
	Note        *string   `json:"note,omitempty"`
	Reason      *string   `json:"reason,omitempty"`
	Meta        *string   `json:"meta,omitempty"`
	Meta2       *string   `json:"meta2,omitempty"`
	Meta3       *string   `json:"meta3,omitempty"`
	Value       *int      `json:"value,omitempty"`
	Rating      *int      `json:"rating,omitempty"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	Deleted     int       `json:"deleted" db:"deleted"`
	TotalCount  int       `json:"total_count,omitempty" db:"total_count"`
}

type OfferResponseUpdate struct {
	State    *string  `json:"state,omitempty"`
	BidPrice *float64 `json:"bid_price,omitempty"`
	Title    *string  `json:"title,omitempty"`
	Note     *string  `json:"note,omitempty"`
	Reason   *string  `json:"reason,omitempty"`
	Meta     *string  `json:"meta,omitempty"`
	Meta2    *string  `json:"meta2,omitempty"`
	Meta3    *string  `json:"meta3,omitempty"`
	Value    *int     `json:"value,omitempty"`
	Rating   *int     `json:"rating,omitempty"`
	Active   *int     `json:"active,omitempty"`
	Deleted  *int     `json:"deleted,omitempty"`
}
