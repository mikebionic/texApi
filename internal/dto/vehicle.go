package dto

type VehicleDetails struct {
	ID                 int           `json:"id"`
	UUID               string        `json:"uuid"`
	CompanyID          int           `json:"company_id"`
	VehicleType        string        `json:"vehicle_type"`
	VehicleBrandID     int           `json:"vehicle_brand_id"`
	VehicleModelID     int           `json:"vehicle_model_id"`
	YearOfIssue        string        `json:"year_of_issue"`
	Mileage            int           `json:"mileage"`
	Numberplate        string        `json:"numberplate"`
	TrailerNumberplate string        `json:"trailer_numberplate"`
	GpsActive          int           `json:"gps_active"`
	Photo1URL          string        `json:"photo1_url"`
	Photo2URL          string        `json:"photo2_url"`
	Photo3URL          string        `json:"photo3_url"`
	Docs1URL           string        `json:"docs1_url"`
	Docs2URL           string        `json:"docs2_url"`
	Docs3URL           string        `json:"docs3_url"`
	ViewCount          int           `json:"view_count"`
	CreatedAt          string        `json:"created_at"`
	UpdatedAt          string        `json:"updated_at"`
	Active             int           `json:"active"`
	Company            *CompanyBasic `json:"company,omitempty"`
	Brand              *VehicleBrand `json:"brand,omitempty"`
	Model              *VehicleModel `json:"model,omitempty"`
}

type VehicleCreate struct {
	CompanyID          int    `json:"company_id" binding:"required"`
	VehicleType        string `json:"vehicle_type" binding:"required"`
	VehicleBrandID     int    `json:"vehicle_brand_id" binding:"required"`
	VehicleModelID     int    `json:"vehicle_model_id" binding:"required"`
	YearOfIssue        string `json:"year_of_issue" binding:"required"`
	Mileage            int    `json:"mileage"`
	Numberplate        string `json:"numberplate" binding:"required"`
	TrailerNumberplate string `json:"trailer_numberplate"`
	GpsActive          int    `json:"gps_active"`
	Photo1URL          string `json:"photo1_url"`
	Photo2URL          string `json:"photo2_url"`
	Photo3URL          string `json:"photo3_url"`
	Docs1URL           string `json:"docs1_url"`
	Docs2URL           string `json:"docs2_url"`
	Docs3URL           string `json:"docs3_url"`
}

type VehicleUpdate struct {
	VehicleType        *string `json:"vehicle_type,omitempty"`
	VehicleBrandID     *int    `json:"vehicle_brand_id,omitempty"`
	VehicleModelID     *int    `json:"vehicle_model_id,omitempty"`
	YearOfIssue        *string `json:"year_of_issue,omitempty"`
	Mileage            *int    `json:"mileage,omitempty"`
	Numberplate        *string `json:"numberplate,omitempty"`
	TrailerNumberplate *string `json:"trailer_numberplate,omitempty"`
	GpsActive          *int    `json:"gps_active,omitempty"`
	Photo1URL          *string `json:"photo1_url,omitempty"`
	Photo2URL          *string `json:"photo2_url,omitempty"`
	Photo3URL          *string `json:"photo3_url,omitempty"`
	Docs1URL           *string `json:"docs1_url,omitempty"`
	Docs2URL           *string `json:"docs2_url,omitempty"`
	Docs3URL           *string `json:"docs3_url,omitempty"`
	Active             *int    `json:"active,omitempty"`
}

type Vehicle struct {
	ID                 int    `json:"id"`
	UUID               string `json:"uuid"`
	VehicleType        string `json:"vehicle_type"`
	VehicleBrandID     int    `json:"vehicle_brand_id"`
	VehicleModelID     int    `json:"vehicle_model_id"`
	YearOfIssue        string `json:"year_of_issue"`
	Mileage            int    `json:"mileage"`
	Numberplate        string `json:"numberplate"`
	TrailerNumberplate string `json:"trailer_numberplate"`
	GpsActive          int    `json:"gps_active"`
	Photo1URL          string `json:"photo1_url"`
	Photo2URL          string `json:"photo2_url"`
	Photo3URL          string `json:"photo3_url"`
	ViewCount          int    `json:"view_count"`
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

type VehicleBasic struct {
	ID          int    `json:"id"`
	VehicleType string `json:"vehicle_type"`
	Numberplate string `json:"numberplate"`
}

//type VehicleBrand struct {
//	ID   int    `json:"id"`
//	Name string `json:"name"`
//}
//
//type VehicleModel struct {
//	ID   int    `json:"id"`
//	Name string `json:"name"`
//}
