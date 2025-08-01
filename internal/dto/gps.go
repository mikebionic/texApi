package dto

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func (p Point) Value() (driver.Value, error) {
	if p.Lat == 0 && p.Lng == 0 {
		return nil, nil
	}
	// PostGIS expects POINT(longitude latitude) format
	return fmt.Sprintf("POINT(%f %f)", p.Lng, p.Lat), nil
}

func (p *Point) Scan(value interface{}) error {
	if value == nil {
		p.Lat = 0
		p.Lng = 0
		return nil
	}

	switch v := value.(type) {
	case string:
		if len(v) > 20 && (v[:2] == "01" || v[:8] == "0101000") {
			// This is WKB hex format, we need to parse it differently
			// For now, return an error asking to use ST_AsText
			return fmt.Errorf("received WKB format, need to use ST_AsText in query")
		}

		// Parse PostGIS text formats
		var lng, lat float64
		if _, err := fmt.Sscanf(v, "POINT(%f %f)", &lng, &lat); err == nil {
			p.Lng = lng
			p.Lat = lat
			return nil
		}
		if _, err := fmt.Sscanf(v, "(%f,%f)", &lng, &lat); err == nil {
			p.Lng = lng
			p.Lat = lat
			return nil
		}
		return fmt.Errorf("invalid point format: %s", v)
	case []byte:
		// Handle binary data
		hexStr := string(v)
		if len(hexStr) > 20 && (hexStr[:2] == "01" || hexStr[:8] == "0101000") {
			return fmt.Errorf("received WKB format, need to use ST_AsText in query")
		}
		return p.Scan(hexStr)
	default:
		return fmt.Errorf("cannot scan %T into Point", value)
	}
}

type Trip struct {
	ID           int64      `json:"id"`
	DriverID     int        `json:"driver_id"`
	VehicleID    int        `json:"vehicle_id"`
	FromAddress  *string    `json:"from_address"`
	ToAddress    *string    `json:"to_address"`
	FromCountry  *string    `json:"from_country"`
	ToCountry    *string    `json:"to_country"`
	StartDate    *time.Time `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	FromLocation *Point     `json:"from_location"`
	ToLocation   *Point     `json:"to_location"`
	DistanceKM   *float64   `json:"distance_km"`
	Status       string     `json:"status"`
	Meta         string     `json:"meta"`
	Meta2        string     `json:"meta2"`
	Meta3        string     `json:"meta3"`
	GPSLogs      string     `json:"gps_logs"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	Deleted      int        `json:"deleted"`
}

type TripOffer struct {
	OfferID int  `json:"offer_id" binding:"required"`
	IsMain  bool `json:"is_main"`
}

type StartTripInput struct {
	DriverID     *int        `json:"driver_id"`
	VehicleID    *int        `json:"vehicle_id"`
	FromAddress  *string     `json:"from_address"`
	ToAddress    *string     `json:"to_address"`
	FromCountry  *string     `json:"from_country"`
	ToCountry    *string     `json:"to_country"`
	StartDate    *time.Time  `json:"start_date"`
	EndDate      *time.Time  `json:"end_date"`
	FromLocation *Point      `json:"from_location"`
	ToLocation   *Point      `json:"to_location"`
	DistanceKM   *float64    `json:"distance_km"`
	Offers       []TripOffer `json:"offers" binding:"required,min=1"`
}

type EndTripInput struct {
	ID        int64 `json:"id" binding:"required"`
	DriverID  int   `json:"driver_id" binding:"required"`
	CompanyID int   `json:"company_id" binding:"required"`
}

type GPSLog struct {
	ID           int64     `json:"id"`
	CompanyID    *int      `json:"company_id"`
	VehicleID    int       `json:"vehicle_id"`
	DriverID     int       `json:"driver_id"`
	OfferID      *int      `json:"offer_id"`
	TripID       *int      `json:"trip_id"`
	BatteryLevel *int      `json:"battery_level"`
	Speed        *float64  `json:"speed"`
	Heading      *float64  `json:"heading"`
	Accuracy     *float64  `json:"accuracy"`
	Coordinates  Point     `json:"coordinates"`
	Status       string    `json:"status"`
	LogDt        time.Time `json:"log_dt"`
	CreatedAt    time.Time `json:"created_at"`
}

type GPSLogInput struct {
	CompanyID    *int      `json:"company_id"`
	VehicleID    int       `json:"vehicle_id" binding:"required"`
	DriverID     int       `json:"driver_id" binding:"required"`
	OfferID      *int      `json:"offer_id"`
	TripID       *int      `json:"trip_id"`
	BatteryLevel *int      `json:"battery_level" binding:"omitempty,min=0,max=100"`
	Speed        *float64  `json:"speed" binding:"omitempty,min=0"`
	Heading      *float64  `json:"heading" binding:"omitempty,min=0,max=359"`
	Accuracy     *float64  `json:"accuracy" binding:"omitempty,min=0"`
	Coordinates  Point     `json:"coordinates" binding:"required"`
	LogDt        time.Time `json:"log_dt" binding:"required"`
}

type GPSLogQuery struct {
	TripID      *int       `form:"trip_id"`
	CompanyID   *int       `form:"company_id"`
	OfferID     *int       `form:"offer_id"`
	DriverID    *int       `form:"driver_id"`
	VehicleID   *int       `form:"vehicle_id"`
	From        *time.Time `form:"from" time_format:"2006-01-02"`
	To          *time.Time `form:"to" time_format:"2006-01-02"`
	TripOfferID *int       `form:"trip_offer_id"`
	Offset      int        `form:"offset" binding:"omitempty,min=0"`
	Limit       int        `form:"limit" binding:"omitempty,min=1,max=1000"`
	OrderBy     *string    `form:"order_by" binding:"omitempty,oneof=id log_dt"`
	OrderDir    *string    `form:"order_dir" binding:"omitempty,oneof=ASC DESC"`
}

type PositionQuery struct {
	TripIDs    []int `form:"trip_ids"`
	CompanyIDs []int `form:"company_ids"`
	OfferIDs   []int `form:"offer_ids"`
	DriverIDs  []int `form:"driver_ids"`
	VehicleIDs []int `form:"vehicle_ids"`
}

type TripQuery struct {
	// Existing fields
	DriverID     *int       `form:"driver_id"`
	VehicleID    *int       `form:"vehicle_id"`
	FromAddress  *string    `form:"from_address"`
	ToAddress    *string    `form:"to_address"`
	FromCountry  *string    `form:"from_country"`
	ToCountry    *string    `form:"to_country"`
	StartDate    *time.Time `form:"start_date" time_format:"2006-01-02"`
	EndDate      *time.Time `form:"end_date" time_format:"2006-01-02"`
	FromLocation *Point     `form:"from_location"`
	ToLocation   *Point     `form:"to_location"`
	DistanceKM   *float64   `form:"distance_km"`
	TripOfferID  *int       `form:"trip_offer_id"`

	Status   *string  `form:"status"`
	StatusIn []string `form:"status_in"`

	StartDateFrom *time.Time `form:"start_date_from" time_format:"2006-01-02"`
	StartDateTo   *time.Time `form:"start_date_to" time_format:"2006-01-02"`
	EndDateFrom   *time.Time `form:"end_date_from" time_format:"2006-01-02"`
	EndDateTo     *time.Time `form:"end_date_to" time_format:"2006-01-02"`

	CreatedAfter  *time.Time `form:"created_after" time_format:"2006-01-02T15:04:05"`
	CreatedBefore *time.Time `form:"created_before" time_format:"2006-01-02T15:04:05"`
	UpdatedAfter  *time.Time `form:"updated_after" time_format:"2006-01-02T15:04:05"`
	UpdatedBefore *time.Time `form:"updated_before" time_format:"2006-01-02T15:04:05"`

	DistanceKMMin *float64 `form:"distance_km_min"`
	DistanceKMMax *float64 `form:"distance_km_max"`

	FromRegion *string `form:"from_region"`
	ToRegion   *string `form:"to_region"`

	NearFromLat *float64 `form:"near_from_lat"`
	NearFromLng *float64 `form:"near_from_lng"`
	FromRadius  *float64 `form:"from_radius"` // in kilometers
	NearToLat   *float64 `form:"near_to_lat"`
	NearToLng   *float64 `form:"near_to_lng"`
	ToRadius    *float64 `form:"to_radius"` // in kilometers

	MetaContains  *string `form:"meta_contains"`
	Meta2Contains *string `form:"meta2_contains"`
	Meta3Contains *string `form:"meta3_contains"`

	HasGPSLogs *bool `form:"has_gps_logs"`

	HasDriver  *bool `form:"has_driver"`
	HasVehicle *bool `form:"has_vehicle"`

	DriverIDs  *string `form:"driver_ids"`
	VehicleIDs *string `form:"vehicle_ids"`
	TripIDs    *string `form:"trip_ids"`

	ExcludeDriverID  *int    `form:"exclude_driver_id"`
	ExcludeVehicleID *int    `form:"exclude_vehicle_id"`
	ExcludeTripIDs   *string `form:"exclude_trip_ids"`

	OfferCompanyID     *int       `form:"offer_company_id"`
	OfferExecCompanyID *int       `form:"offer_exec_company_id"`
	OfferDriverID      *int       `form:"offer_driver_id"`
	OfferVehicleID     *int       `form:"offer_vehicle_id"`
	OfferFromCountryID *int       `form:"offer_from_country_id"`
	OfferToCountryID   *int       `form:"offer_to_country_id"`
	OfferFromAddress   *string    `form:"offer_from_address"`
	OfferToAddress     *string    `form:"offer_to_address"`
	OfferFromCountry   *string    `form:"offer_from_country"`
	OfferToCountry     *string    `form:"offer_to_country"`
	OfferState         *string    `form:"offer_state"`
	OfferRole          *string    `form:"offer_role"`
	OfferValidityStart *time.Time `form:"offer_validity_start" time_format:"2006-01-02"`
	OfferValidityEnd   *time.Time `form:"offer_validity_end" time_format:"2006-01-02"`
	OfferDeliveryStart *time.Time `form:"offer_delivery_start" time_format:"2006-01-02"`
	OfferDeliveryEnd   *time.Time `form:"offer_delivery_end" time_format:"2006-01-02"`

	Search *string `form:"search"`

	OfferCostPerKmMin  *float64 `form:"offer_cost_per_km_min"`
	OfferCostPerKmMax  *float64 `form:"offer_cost_per_km_max"`
	OfferPriceMin      *float64 `form:"offer_price_min"`
	OfferPriceMax      *float64 `form:"offer_price_max"`
	OfferTotalPriceMin *float64 `form:"offer_total_price_min"`
	OfferTotalPriceMax *float64 `form:"offer_total_price_max"`

	Offset   int     `form:"offset" binding:"omitempty,min=0"`
	Limit    int     `form:"limit" binding:"omitempty,min=1,max=1000"`
	OrderBy  *string `form:"order_by" binding:"omitempty,oneof=id start_date end_date distance_km created_at updated_at"`
	OrderDir *string `form:"order_dir" binding:"omitempty,oneof=ASC DESC"`
}

type TripDetailed struct {
	ID           int64            `json:"id"`
	DriverID     int              `json:"driver_id"`
	VehicleID    int              `json:"vehicle_id"`
	FromAddress  *string          `json:"from_address"`
	ToAddress    *string          `json:"to_address"`
	FromCountry  *string          `json:"from_country"`
	ToCountry    *string          `json:"to_country"`
	StartDate    *time.Time       `json:"start_date"`
	EndDate      *time.Time       `json:"end_date"`
	FromLocation *Point           `json:"from_location"`
	ToLocation   *Point           `json:"to_location"`
	DistanceKM   *float64         `json:"distance_km"`
	Status       string           `json:"status"`
	Meta         string           `json:"meta"`
	Meta2        string           `json:"meta2"`
	Meta3        string           `json:"meta3"`
	GPSLogs      string           `json:"gps_logs"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	Deleted      int              `json:"deleted"`
	TotalCount   int              `json:"total_count"`
	Driver       *json.RawMessage `json:"driver"`
	Vehicle      *json.RawMessage `json:"vehicle"`
	Offers       *json.RawMessage `json:"offers"`
}
