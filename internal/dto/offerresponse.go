package dto

import "time"

type OfferResponseDetails struct {
	ID          int           `json:"id" db:"id"`
	UUID        string        `json:"uuid" db:"uuid"`
	CompanyID   int           `json:"company_id" db:"company_id"`
	OfferID     int           `json:"offer_id" db:"offer_id"`
	ToCompanyID int           `json:"to_company_id" db:"to_company_id"`
	State       string        `json:"state" db:"state"`
	BidPrice    float64       `json:"bid_price" db:"bid_price"`
	Title       string        `json:"title" db:"title"`
	Note        string        `json:"note" db:"note"`
	Reason      string        `json:"reason" db:"reason"`
	Meta        string        `json:"meta" db:"meta"`
	Meta2       string        `json:"meta2" db:"meta2"`
	Meta3       string        `json:"meta3" db:"meta3"`
	Value       int           `json:"value" db:"value"`
	Rating      int           `json:"rating" db:"rating"`
	CreatedAt   time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" db:"updated_at"`
	Deleted     int           `json:"deleted" db:"deleted"`
	TotalCount  int           `json:"total_count" db:"total_count"`
	Company     CompanyCreate `json:"company" db:"company"`
	ToCompany   CompanyCreate `json:"to_company" db:"to_company"`
	Offer       Offer         `json:"offer" db:"offer"`
}

type OfferResponseCreate struct {
	CompanyID   int     `json:"company_id"`
	OfferID     int     `json:"offer_id"`
	ToCompanyID int     `json:"to_company_id"`
	State       string  `json:"state"`
	BidPrice    float64 `json:"bid_price"`
	Title       string  `json:"title"`
	Note        string  `json:"note"`
	Reason      string  `json:"reason,omitempty"`
	Meta        string  `json:"meta,omitempty"`
	Meta2       string  `json:"meta2,omitempty"`
	Meta3       string  `json:"meta3,omitempty"`
	Value       int     `json:"value,omitempty"`
	Rating      int     `json:"rating,omitempty"`
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
