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
	VehicleTypeID    int       `json:"vehicle_type_id"`
	CargoID          int       `json:"cargo_id"`
	PackagingTypeID  int       `json:"packaging_type_id"`
	OfferState       string    `json:"offer_state"`
	OfferRole        string    `json:"offer_role"`
	CostPerKm        float64   `json:"cost_per_km"`
	Currency         string    `json:"currency"`
	FromCountryID    int       `json:"from_country_id"`
	FromCityID       int       `json:"from_city_id"`
	ToCountryID      int       `json:"to_country_id"`
	ToCityID         int       `json:"to_city_id"`
	Distance         int       `json:"distance"`
	FromCountry      string    `json:"from_country"`
	FromRegion       string    `json:"from_region"`
	ToCountry        string    `json:"to_country"`
	ToRegion         string    `json:"to_region"`
	FromAddress      string    `json:"from_address"`
	ToAddress        string    `json:"to_address"`
	MapURL           string    `json:"map_url"`
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
	PaymentTerm      string    `json:"payment_term"`
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

type CompanyWithStats struct {
	CompanyCreate
	DriversCount      int `json:"drivers_count"`
	ActiveOffersCount int `json:"offers_count"`
}

type OfferDetailedResponse struct {
	Offer
	Company        *CompanyWithStats      `json:"company,omitempty"`
	ExecCompany    *CompanyWithStats      `json:"exec_company,omitempty"`
	AssignedDriver *DriverCreate          `json:"assigned_driver,omitempty"`
	Vehicle        *VehicleCreate         `json:"vehicle,omitempty"`
	VehicleType    *VehicleType           `json:"vehicle_type,omitempty"`
	Cargo          *Cargo                 `json:"cargo,omitempty"`
	PackagingType  *PackagingTypeResponse `json:"packaging_type,omitempty"`
}

type OfferDetails struct {
	Offer
	Company         *CompanyBasic `json:"company,omitempty"`
	AssignedDriver  *DriverShort  `json:"assigned_driver,omitempty"`
	AssignedVehicle *VehicleBasic `json:"assigned_vehicle,omitempty"`
	Cargo           *Cargo        `json:"cargo,omitempty"`
}

// OfferUpdate represents the structure for updating an existing offer
type OfferUpdate struct {
	ID               *int     `json:"id,omitempty"`
	UserID           *int     `json:"user_id,omitempty"`
	CompanyID        *int     `json:"company_id,omitempty"`
	ExecCompanyID    *int     `json:"exec_company_id,omitempty"`
	DriverID         *int     `json:"driver_id,omitempty"`
	VehicleID        *int     `json:"vehicle_id,omitempty"`
	VehicleTypeID    *int     `json:"vehicle_type_id,omitempty"`
	CargoID          *int     `json:"cargo_id,omitempty"`
	PackagingTypeID  *int     `json:"packaging_type_id,omitempty"`
	OfferState       *string  `json:"offer_state,omitempty"`
	OfferRole        *string  `json:"offer_role,omitempty"`
	CostPerKm        *float64 `json:"cost_per_nkm,omitempty"`
	Currency         *string  `json:"currency,omitempty"`
	FromCountryID    *int     `json:"from_country_id,omitempty"`
	FromCityID       *int     `json:"from_city_id,omitempty"`
	ToCountryID      *int     `json:"to_country_id,omitempty"`
	ToCityID         *int     `json:"to_city_id,omitempty"`
	Distance         *int     `json:"distance,omitempty"`
	FromCountry      *string  `json:"from_country,omitempty"`
	FromRegion       *string  `json:"from_region,omitempty"`
	ToCountry        *string  `json:"to_country,omitempty"`
	ToRegion         *string  `json:"to_region,omitempty"`
	FromAddress      *string  `json:"from_address,omitempty"`
	ToAddress        *string  `json:"to_address,omitempty"`
	MapURL           *string  `json:"map_url,omitempty"`
	SenderContact    *string  `json:"sender_contact,omitempty"`
	RecipientContact *string  `json:"recipient_contact,omitempty"`
	DeliverContact   *string  `json:"deliver_contact,omitempty"`
	ViewCount        *int     `json:"view_count,omitempty"`
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
	PaymentTerm      *string  `json:"payment_term,omitempty"`
	Meta             *string  `json:"meta,omitempty"`
	Meta2            *string  `json:"meta2,omitempty"`
	Meta3            *string  `json:"meta3,omitempty"`
	Featured         *int     `json:"featured,omitempty"`
	Partner          *int     `json:"partner,omitempty"`
	CreatedAt        *string  `json:"created_at,omitempty"`
	UpdatedAt        *string  `json:"updated_at,omitempty"`
	Active           *int     `json:"active,omitempty"`
	Deleted          *int     `json:"deleted,omitempty"`
}
