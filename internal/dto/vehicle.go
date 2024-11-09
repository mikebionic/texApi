package dto

type VehicleCreate struct {
	ID                 int    `json:"id"`
	CompanyID          int    `json:"company_id"`
	VehicleType        string `json:"vehicle_type"`
	Brand              string `json:"brand"`
	VehicleModel       string `json:"vehicle_model"`
	YearOfIssue        string `json:"year_of_issue"`
	Numberplate        string `json:"numberplate"`
	Mileage            string `json:"mileage"`
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

type VehicleBrand struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Country     string `json:"country"`
	FoundedYear int    `json:"founded_year"`
	Deleted     int    `json:"deleted"`
}

type VehicleType struct {
	ID          int    `json:"id"`
	TypeName    string `json:"type_name"`
	Description string `json:"description"`
	Deleted     int    `json:"deleted"`
}

type VehicleModel struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Year          int    `json:"year"`
	Brand         string `json:"brand"`
	VehicleTypeID int    `json:"vehicle_type_id"`
	VehicleType   string `json:"vehicle_type"`
	Feature       string `json:"feature"`
	Deleted       int    `json:"deleted"`
}

type VehicleBrandUpdate struct {
	Name        *string `json:"name,omitempty"`
	Country     *string `json:"country,omitempty"`
	FoundedYear *int    `json:"founded_year,omitempty"`
}

type VehicleTypeUpdate struct {
	TypeName    *string `json:"type_name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type VehicleModelUpdate struct {
	Name          *string `json:"name,omitempty"`
	Year          *int    `json:"year,omitempty"`
	Brand         *string `json:"brand,omitempty"`
	VehicleTypeID *int    `json:"vehicle_type_id,omitempty"`
	Feature       *string `json:"feature,omitempty"`
}
