package dto

import (
	"time"
)

type PriceQuote struct {
	ID                 int       `json:"id"`
	UUID               string    `json:"uuid"`
	TransportType      string    `json:"transport_type"`
	SubType            string    `json:"sub_type"`
	UserID             int       `json:"user_id"`
	CompanyID          int       `json:"company_id"`
	ExecCompanyID      int       `json:"exec_company_id"`
	VehicleTypeID      int       `json:"vehicle_type_id"`
	PackagingTypeID    int       `json:"packaging_type_id"`
	CostPerKm          float64   `json:"cost_per_km"`
	Currency           string    `json:"currency"`
	FromCountryID      int       `json:"from_country_id"`
	FromCityID         int       `json:"from_city_id"`
	ToCountryID        int       `json:"to_country_id"`
	ToCityID           int       `json:"to_city_id"`
	Distance           int       `json:"distance"`
	FromCountry        string    `json:"from_country"`
	FromRegion         string    `json:"from_region"`
	ToCountry          string    `json:"to_country"`
	ToRegion           string    `json:"to_region"`
	FromAddress        string    `json:"from_address"`
	ToAddress          string    `json:"to_address"`
	Tax                int       `json:"tax"`
	TaxPrice           float64   `json:"tax_price"`
	Trade              int       `json:"trade"`
	Discount           int       `json:"discount"`
	PaymentMethod      string    `json:"payment_method"`
	PaymentTerm        string    `json:"payment_term"`
	DistanceKm         int       `json:"distance_km"`
	AveragePrice       float64   `json:"average_price"`
	MinPrice           float64   `json:"min_price"`
	MaxPrice           float64   `json:"max_price"`
	PriceUnit          string    `json:"price_unit"`
	MinVolume          float64   `json:"min_volume"`
	MaxVolume          float64   `json:"max_volume"`
	ValidityStart      time.Time `json:"validity_start"`
	ValidityEnd        time.Time `json:"validity_end"`
	FuelIncluded       bool      `json:"fuel_included"`
	CustomsIncluded    bool      `json:"customs_included"`
	InsuranceIncluded  bool      `json:"insurance_included"`
	FuelInfo           string    `json:"fuel_info"`
	CustomsInfo        string    `json:"customs_info"`
	InsuranceInfo      string    `json:"insurance_info"`
	Terms              string    `json:"terms"`
	SurchargeInfo      string    `json:"surcharge_info"`
	IsPromotional      bool      `json:"is_promotional"`
	IsDynamic          bool      `json:"is_dynamic"`
	DataSource         string    `json:"data_source"`
	UpdatedFromOfferID int       `json:"updated_from_offer_id"`
	SampleSize         int       `json:"sample_size"`
	Notes              string    `json:"notes"`
	InternalNote       string    `json:"internal_note"`
	Meta               string    `json:"meta"`
	Meta2              string    `json:"meta2"`
	Meta3              string    `json:"meta3"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	Active             int       `json:"active"`
	Deleted            int       `json:"deleted"`
}

type CreatePriceQuoteRequest struct {
	TransportType      string    `json:"transport_type" binding:"required"`
	SubType            string    `json:"sub_type"`
	UserID             int       `json:"user_id"`
	CompanyID          int       `json:"company_id"`
	ExecCompanyID      int       `json:"exec_company_id"`
	VehicleTypeID      int       `json:"vehicle_type_id"`
	PackagingTypeID    int       `json:"packaging_type_id"`
	CostPerKm          float64   `json:"cost_per_km"`
	Currency           string    `json:"currency"`
	FromCountryID      int       `json:"from_country_id"`
	FromCityID         int       `json:"from_city_id"`
	ToCountryID        int       `json:"to_country_id"`
	ToCityID           int       `json:"to_city_id"`
	Distance           int       `json:"distance"`
	FromCountry        string    `json:"from_country"`
	FromRegion         string    `json:"from_region"`
	ToCountry          string    `json:"to_country"`
	ToRegion           string    `json:"to_region"`
	FromAddress        string    `json:"from_address"`
	ToAddress          string    `json:"to_address"`
	Tax                int       `json:"tax"`
	TaxPrice           float64   `json:"tax_price"`
	Trade              int       `json:"trade"`
	Discount           int       `json:"discount"`
	PaymentMethod      string    `json:"payment_method"`
	PaymentTerm        string    `json:"payment_term"`
	DistanceKm         int       `json:"distance_km"`
	AveragePrice       float64   `json:"average_price" binding:"required"`
	MinPrice           float64   `json:"min_price"`
	MaxPrice           float64   `json:"max_price"`
	PriceUnit          string    `json:"price_unit" binding:"required"`
	MinVolume          float64   `json:"min_volume"`
	MaxVolume          float64   `json:"max_volume"`
	ValidityStart      time.Time `json:"validity_start"`
	ValidityEnd        time.Time `json:"validity_end"`
	FuelIncluded       bool      `json:"fuel_included"`
	CustomsIncluded    bool      `json:"customs_included"`
	InsuranceIncluded  bool      `json:"insurance_included"`
	FuelInfo           string    `json:"fuel_info"`
	CustomsInfo        string    `json:"customs_info"`
	InsuranceInfo      string    `json:"insurance_info"`
	Terms              string    `json:"terms"`
	SurchargeInfo      string    `json:"surcharge_info"`
	IsPromotional      bool      `json:"is_promotional"`
	IsDynamic          bool      `json:"is_dynamic"`
	DataSource         string    `json:"data_source"`
	UpdatedFromOfferID int       `json:"updated_from_offer_id"`
	SampleSize         int       `json:"sample_size"`
	Notes              string    `json:"notes"`
	InternalNote       string    `json:"internal_note"`
	Meta               string    `json:"meta"`
	Meta2              string    `json:"meta2"`
	Meta3              string    `json:"meta3"`
}

type UpdatePriceQuoteRequest struct {
	TransportType      *string    `json:"transport_type"`
	SubType            *string    `json:"sub_type"`
	UserID             *int       `json:"user_id"`
	CompanyID          *int       `json:"company_id"`
	ExecCompanyID      *int       `json:"exec_company_id"`
	VehicleTypeID      *int       `json:"vehicle_type_id"`
	PackagingTypeID    *int       `json:"packaging_type_id"`
	CostPerKm          *float64   `json:"cost_per_km"`
	Currency           *string    `json:"currency"`
	FromCountryID      *int       `json:"from_country_id"`
	FromCityID         *int       `json:"from_city_id"`
	ToCountryID        *int       `json:"to_country_id"`
	ToCityID           *int       `json:"to_city_id"`
	Distance           *int       `json:"distance"`
	FromCountry        *string    `json:"from_country"`
	FromRegion         *string    `json:"from_region"`
	ToCountry          *string    `json:"to_country"`
	ToRegion           *string    `json:"to_region"`
	FromAddress        *string    `json:"from_address"`
	ToAddress          *string    `json:"to_address"`
	Tax                *int       `json:"tax"`
	TaxPrice           *float64   `json:"tax_price"`
	Trade              *int       `json:"trade"`
	Discount           *int       `json:"discount"`
	PaymentMethod      *string    `json:"payment_method"`
	PaymentTerm        *string    `json:"payment_term"`
	DistanceKm         *int       `json:"distance_km"`
	AveragePrice       *float64   `json:"average_price"`
	MinPrice           *float64   `json:"min_price"`
	MaxPrice           *float64   `json:"max_price"`
	PriceUnit          *string    `json:"price_unit"`
	MinVolume          *float64   `json:"min_volume"`
	MaxVolume          *float64   `json:"max_volume"`
	ValidityStart      *time.Time `json:"validity_start"`
	ValidityEnd        *time.Time `json:"validity_end"`
	FuelIncluded       *bool      `json:"fuel_included"`
	CustomsIncluded    *bool      `json:"customs_included"`
	InsuranceIncluded  *bool      `json:"insurance_included"`
	FuelInfo           *string    `json:"fuel_info"`
	CustomsInfo        *string    `json:"customs_info"`
	InsuranceInfo      *string    `json:"insurance_info"`
	Terms              *string    `json:"terms"`
	SurchargeInfo      *string    `json:"surcharge_info"`
	IsPromotional      *bool      `json:"is_promotional"`
	IsDynamic          *bool      `json:"is_dynamic"`
	DataSource         *string    `json:"data_source"`
	UpdatedFromOfferID *int       `json:"updated_from_offer_id"`
	SampleSize         *int       `json:"sample_size"`
	Notes              *string    `json:"notes"`
	InternalNote       *string    `json:"internal_note"`
	Meta               *string    `json:"meta"`
	Meta2              *string    `json:"meta2"`
	Meta3              *string    `json:"meta3"`
	Active             *int       `json:"active"`
}

type PriceQuoteFilters struct {
	TransportType   string  `form:"transport_type"`
	SubType         string  `form:"sub_type"`
	FromCountry     string  `form:"from_country"`
	ToCountry       string  `form:"to_country"`
	FromRegion      string  `form:"from_region"`
	ToRegion        string  `form:"to_region"`
	Currency        string  `form:"currency"`
	PriceUnit       string  `form:"price_unit"`
	MinPrice        float64 `form:"min_price"`
	MaxPrice        float64 `form:"max_price"`
	FuelIncluded    *bool   `form:"fuel_included"`
	CustomsIncluded *bool   `form:"customs_included"`
	IsPromotional   *bool   `form:"is_promotional"`
	IsDynamic       *bool   `form:"is_dynamic"`
	DataSource      string  `form:"data_source"`
	UserID          int     `form:"user_id"`
	CompanyID       int     `form:"company_id"`
	VehicleTypeID   int     `form:"vehicle_type_id"`
	Active          *int    `form:"active"`
	Page            int     `form:"page"`
	PerPage         int     `form:"per_page"`
	SortBy          string  `form:"sort_by"`
	SortOrder       string  `form:"sort_order"`
	PaymentMethod   string  `form:"payment_method"`
}

type PriceQuoteAnalysisFilters struct {
	TransportType   string  `form:"transport_type"`
	SubType         string  `form:"sub_type"`
	FromCountryID   int     `form:"from_country_id"`
	FromCityID      int     `form:"from_city_id"`
	ToCountryID     int     `form:"to_country_id"`
	ToCityID        int     `form:"to_city_id"`
	FromCountry     string  `form:"from_country"`
	ToCountry       string  `form:"to_country"`
	FromRegion      string  `form:"from_region"`
	ToRegion        string  `form:"to_region"`
	Currency        string  `form:"currency"`
	VehicleTypeID   int     `form:"vehicle_type_id"`
	PackagingTypeID int     `form:"packaging_type_id"`
	PaymentMethod   string  `form:"payment_method"`
	Distance        int     `form:"distance"`
	DistanceKm      int     `form:"distance_km"`
	MinVolume       float64 `form:"min_volume"`
	MaxVolume       float64 `form:"max_volume"`
	FuelIncluded    *bool   `form:"fuel_included"`
	CustomsIncluded *bool   `form:"customs_included"`
	// Matching tolerance (how strict the matching should be)
	MatchStrict bool `form:"match_strict"` // if true, requires exact matches
}

type PriceQuoteAnalysisResponse struct {
	PriceQuote
	AnalysisInfo struct {
		FoundFromOffers    bool     `json:"found_from_offers"`
		OfferCount         int      `json:"offer_count"`
		PriceQuoteCount    int      `json:"price_quote_count"`
		OfferMinPrice      float64  `json:"offer_min_price,omitempty"`
		OfferMaxPrice      float64  `json:"offer_max_price,omitempty"`
		OfferAvgCostPerKm  float64  `json:"offer_avg_cost_per_km,omitempty"`
		PriceQuoteMinPrice float64  `json:"price_quote_min_price,omitempty"`
		PriceQuoteMaxPrice float64  `json:"price_quote_max_price,omitempty"`
		PriceQuoteAvgPrice float64  `json:"price_quote_avg_price,omitempty"`
		MatchingCriteria   []string `json:"matching_criteria"`
	} `json:"analysis_info"`
}
