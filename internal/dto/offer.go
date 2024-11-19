package dto

import "time"

type Offer struct {
	ID               int       `json:"id"`
	UUID             string    `json:"uuid"`
	UserID           int       `json:"user_id"`
	CompanyID        int       `json:"company_id"`
	ExecCompanyID    int       `json:"exec_company_id"`
	DriverID         int       `json:"driver_id"`
	VehicleID        int       `json:"vehicle_id"`
	CargoID          int       `json:"cargo_id"`
	OfferState       string    `json:"offer_state"`
	OfferRole        string    `json:"offer_role"`
	CostPerKm        float64   `json:"cost_per_km"`
	Currency         string    `json:"currency"`
	FromCountryID    int       `json:"from_country_id"`
	FromCityID       int       `json:"from_city_id"`
	ToCountryID      int       `json:"to_country_id"`
	ToCityID         int       `json:"to_city_id"`
	FromCountry      string    `json:"from_country"`
	FromRegion       string    `json:"from_region"`
	ToCountry        string    `json:"to_country"`
	ToRegion         string    `json:"to_region"`
	FromAddress      string    `json:"from_address"`
	ToAddress        string    `json:"to_address"`
	SenderContact    string    `json:"sender_contact"`
	RecipientContact string    `json:"recipient_contact"`
	DeliverContact   string    `json:"deliver_contact"`
	ViewCount        int       `json:"view_count"`
	ValidityStart    time.Time `json:"validity_start"`
	ValidityEnd      time.Time `json:"validity_end"`
	DeliveryStart    time.Time `json:"delivery_start"`
	DeliveryEnd      time.Time `json:"delivery_end"`
	Note             string    `json:"note"`
	Tax              int       `json:"tax"`
	TaxPrice         float64   `json:"tax_price"`
	Trade            int       `json:"trade"`
	Discount         int       `json:"discount"`
	PaymentMethod    string    `json:"payment_method"`
	Meta             string    `json:"meta"`
	Meta2            string    `json:"meta2"`
	Meta3            string    `json:"meta3"`
	Featured         int       `json:"featured"`
	Partner          int       `json:"partner"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	Active           int       `json:"active"`
	Deleted          int       `json:"deleted"`
	TotalCount       int       `json:"total_count"`
}

type OfferDetails struct {
	Offer
	Company         *CompanyBasic `json:"company,omitempty"`
	AssignedDriver  *Driver       `json:"assigned_driver,omitempty"`
	AssignedVehicle *VehicleBasic `json:"assigned_vehicle,omitempty"`
	Cargo           *Cargo        `json:"cargo,omitempty"`
}

type OfferUpdate struct {
	DriverID         *int     `json:"driver_id,omitempty"`
	CompanyID        *int     `json:"company_id,omitempty"`
	ExecCompanyID    *int     `json:"exec_company_id,omitempty"`
	VehicleID        *int     `json:"vehicle_id,omitempty"`
	CargoID          *int     `json:"cargo_id,omitempty"`
	OfferState       *string  `json:"offer_state,omitempty"`
	OfferRole        *string  `json:"offer_role,omitempty"`
	CostPerKm        *float64 `json:"cost_per_km,omitempty"`
	Currency         *string  `json:"currency,omitempty"`
	FromCountryID    *int     `json:"from_country_id,omitempty"`
	FromCityID       *int     `json:"from_city_id,omitempty"`
	ToCountryID      *int     `json:"to_country_id,omitempty"`
	ToCityID         *int     `json:"to_city_id,omitempty"`
	FromCountry      *string  `json:"from_country,omitempty"`
	FromRegion       *string  `json:"from_region,omitempty"`
	ToCountry        *string  `json:"to_country,omitempty"`
	ToRegion         *string  `json:"to_region,omitempty"`
	FromAddress      *string  `json:"from_address,omitempty"`
	ToAddress        *string  `json:"to_address,omitempty"`
	SenderContact    *string  `json:"sender_contact,omitempty"`
	RecipientContact *string  `json:"recipient_contact,omitempty"`
	DeliverContact   *string  `json:"deliver_contact,omitempty"`
	ValidityStart    *string  `json:"validity_start,omitempty"`
	ValidityEnd      *string  `json:"validity_end,omitempty"`
	DeliveryStart    *string  `json:"delivery_start,omitempty"`
	DeliveryEnd      *string  `json:"delivery_end,omitempty"`
	Note             *string  `json:"note,omitempty"`
	Tax              *int     `json:"tax,omitempty"`
	TaxPrice         *float64 `json:"tax_price,omitempty"`
	Trade            *int     `json:"trade,omitempty"`
	Discount         *int     `json:"discount,omitempty"`
	PaymentMethod    *string  `json:"payment_method,omitempty"`
	Meta             *string  `json:"meta,omitempty"`
	Meta2            *string  `json:"meta2,omitempty"`
	Meta3            *string  `json:"meta3,omitempty"`
	Active           *int     `json:"active,omitempty"`
	Deleted          *int     `json:"deleted,omitempty"`
}
