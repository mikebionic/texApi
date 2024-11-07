package dto

type RequestCreate struct {
	ID            int     `json:"id"`
	UserID        int     `json:"user_id"`
	CompanyID     int     `json:"company_id"`
	DriverID      int     `json:"driver_id"`
	VehicleID     int     `json:"vehicle_id"`
	CostPerKM     float64 `json:"cost_per_km"`
	FromCountry   string  `json:"from_country"`
	FromRegion    string  `json:"from_region"`
	ToCountry     string  `json:"to_country"`
	ToRegion      string  `json:"to_region"`
	ViewCount     string  `json:"view_count"`
	ValidityStart string  `json:"validity_start"`
	ValidityEnd   string  `json:"validity_end"`
	Note          string  `json:"note"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	Deleted       int     `json:"deleted"`
}

type RequestUpdate struct {
	DriverID      *int     `json:"driver_id,omitempty"`
	VehicleID     *int     `json:"vehicle_id,omitempty"`
	CostPerKM     *float64 `json:"cost_per_km,omitempty"`
	FromCountry   *string  `json:"from_country,omitempty"`
	FromRegion    *string  `json:"from_region,omitempty"`
	ToCountry     *string  `json:"to_country,omitempty"`
	ToRegion      *string  `json:"to_region,omitempty"`
	ValidityStart *string  `json:"validity_start,omitempty"`
	ValidityEnd   *string  `json:"validity_end,omitempty"`
	Note          *string  `json:"note,omitempty"`
}
