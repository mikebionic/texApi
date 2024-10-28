package dto

type VehicleCreate struct {
	ID                 int    `json:"id"`
	CompanyID          int    `json:"company_id"`
	VehicleType        string `json:"vehicle_type"`
	Brand              string `json:"brand"`
	VehicleModel       string `json:"vehicle_model"`
	YearOfIssue        string `json:"year_of_issue"`
	Numberplate        string `json:"numberplate"`
	TrailerNumberplate string `json:"trailer_numberplate"`
	GPSActive          int    `json:"gps_active"`
	Photo1URL          string `json:"photo1_url"`
	Photo2URL          string `json:"photo2_url"`
	Photo3URL          string `json:"photo3_url"`
	Docs1URL           string `json:"docs1_url"`
	Docs2URL           string `json:"docs2_url"`
	Docs3URL           string `json:"docs3_url"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
	Active             int    `json:"active"`
	Deleted            int    `json:"deleted"`
}

type VehicleUpdate struct {
	VehicleType        *string `json:"vehicle_type,omitempty"`
	Brand              *string `json:"brand,omitempty"`
	VehicleModel       *string `json:"vehicle_model,omitempty"`
	YearOfIssue        *string `json:"year_of_issue,omitempty"`
	Numberplate        *string `json:"numberplate,omitempty"`
	TrailerNumberplate *string `json:"trailer_numberplate,omitempty"`
	GPSActive          *int    `json:"gps_active,omitempty"`
	Photo1URL          *string `json:"photo1_url,omitempty"`
	Photo2URL          *string `json:"photo2_url,omitempty"`
	Photo3URL          *string `json:"photo3_url,omitempty"`
	Docs1URL           *string `json:"docs1_url,omitempty"`
	Docs2URL           *string `json:"docs2_url,omitempty"`
	Docs3URL           *string `json:"docs3_url,omitempty"`
}
