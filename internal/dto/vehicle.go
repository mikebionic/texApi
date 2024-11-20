package dto

import "time"

type VehicleDetails struct {
	VehicleCreate
	Company *CompanyBasic `json:"company,omitempty"`
	Brand   *VehicleBrand `json:"brand,omitempty"`
	Model   *VehicleModel `json:"model,omitempty"`
}

type VehicleCreate struct {
	ID                 int       `json:"id"`
	UUID               string    `json:"uuid"`
	CompanyID          int       `json:"company_id"`
	VehicleTypeID      int       `json:"vehicle_type_id"`
	VehicleBrandID     int       `json:"vehicle_brand_id"`
	VehicleModelID     int       `json:"vehicle_model_id"`
	YearOfIssue        string    `json:"year_of_issue"`
	Mileage            int       `json:"mileage"`
	Numberplate        string    `json:"numberplate"`
	TrailerNumberplate string    `json:"trailer_numberplate"`
	Gps                int       `json:"gps"`
	Photo1URL          string    `json:"photo1_url"`
	Photo2URL          string    `json:"photo2_url"`
	Photo3URL          string    `json:"photo3_url"`
	Docs1URL           string    `json:"docs1_url"`
	Docs2URL           string    `json:"docs2_url"`
	Docs3URL           string    `json:"docs3_url"`
	ViewCount          int       `json:"view_count"`
	Meta               string    `json:"meta"`
	Meta2              string    `json:"meta2"`
	Meta3              string    `json:"meta3"`
	Available          int       `json:"available"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	Active             int       `json:"active"`
	Deleted            int       `json:"deleted"`
	TotalCount         int       `json:"total_count"`
}

type VehicleUpdate struct {
	CompanyID          *int    `json:"company_id,omitempty"`
	VehicleTypeID      *int    `json:"vehicle_type_id,omitempty"`
	VehicleBrandID     *int    `json:"vehicle_brand_id,omitempty"`
	VehicleModelID     *int    `json:"vehicle_model_id,omitempty"`
	YearOfIssue        *string `json:"year_of_issue,omitempty"`
	Mileage            *int    `json:"mileage,omitempty"`
	Numberplate        *string `json:"numberplate,omitempty"`
	TrailerNumberplate *string `json:"trailer_numberplate,omitempty"`
	Gps                *int    `json:"gps,omitempty"`
	Photo1URL          *string `json:"photo1_url,omitempty"`
	Photo2URL          *string `json:"photo2_url,omitempty"`
	Photo3URL          *string `json:"photo3_url,omitempty"`
	Docs1URL           *string `json:"docs1_url,omitempty"`
	Docs2URL           *string `json:"docs2_url,omitempty"`
	Docs3URL           *string `json:"docs3_url,omitempty"`
	ViewCount          *int    `json:"view_count"`
	Meta               *string `json:"meta"`
	Meta2              *string `json:"meta2"`
	Meta3              *string `json:"meta3"`
	Available          *int    `json:"available"`
	Active             *int    `json:"active,omitempty"`
	Deleted            *int    `json:"deleted,omitempty"`
}

type VehicleShort struct {
	ID                 int    `json:"id"`
	UUID               string `json:"uuid"`
	VehicleTypeID      int    `json:"vehicle_type_id"`
	VehicleBrandID     int    `json:"vehicle_brand_id"`
	VehicleModelID     int    `json:"vehicle_model_id"`
	YearOfIssue        string `json:"year_of_issue"`
	Mileage            int    `json:"mileage"`
	Numberplate        string `json:"numberplate"`
	TrailerNumberplate string `json:"trailer_numberplate"`
	Gps                int    `json:"gps"`
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
	ID      int    `json:"id"`
	TitleEn string `json:"title_en"`
	DescEn  string `json:"desc_en"`
	TitleRu string `json:"title_ru"`
	DescRu  string `json:"desc_ru"`
	TitleTk string `json:"title_tk"`
	DescTk  string `json:"desc_tk"`
	TitleDe string `json:"title_de"`
	DescDe  string `json:"desc_de"`
	TitleAr string `json:"title_ar"`
	DescAr  string `json:"desc_ar"`
	TitleEs string `json:"title_es"`
	DescEs  string `json:"desc_es"`
	TitleFr string `json:"title_fr"`
	DescFr  string `json:"desc_fr"`
	TitleZh string `json:"title_zh"`
	DescZh  string `json:"desc_zh"`
	TitleJa string `json:"title_ja"`
	DescJa  string `json:"desc_ja"`
	Deleted int    `json:"deleted"`
}

type VehicleTypeUpdate struct {
	TitleEn *string `json:"title_en,omitempty"`
	DescEn  *string `json:"desc_en,omitempty"`
	TitleRu *string `json:"title_ru,omitempty"`
	DescRu  *string `json:"desc_ru,omitempty"`
	TitleTk *string `json:"title_tk,omitempty"`
	DescTk  *string `json:"desc_tk,omitempty"`
	TitleDe *string `json:"title_de,omitempty"`
	DescDe  *string `json:"desc_de,omitempty"`
	TitleAr *string `json:"title_ar,omitempty"`
	DescAr  *string `json:"desc_ar,omitempty"`
	TitleEs *string `json:"title_es,omitempty"`
	DescEs  *string `json:"desc_es,omitempty"`
	TitleFr *string `json:"title_fr,omitempty"`
	DescFr  *string `json:"desc_fr,omitempty"`
	TitleZh *string `json:"title_zh,omitempty"`
	DescZh  *string `json:"desc_zh,omitempty"`
	TitleJa *string `json:"title_ja,omitempty"`
	DescJa  *string `json:"desc_ja,omitempty"`
}

type VehicleModel struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Year           int    `json:"year"`
	Brand          string `json:"brand"`
	VehicleBrandID int    `json:"vehicle_brand_id"`
	VehicleBrand   string `json:"vehicle_brand"`
	VehicleTypeID  int    `json:"vehicle_type_id"`
	VehicleType    string `json:"vehicle_type"`
	Feature        string `json:"feature"`
	Deleted        int    `json:"deleted"`
}

type VehicleBrandUpdate struct {
	Name        *string `json:"name,omitempty"`
	Country     *string `json:"country,omitempty"`
	FoundedYear *int    `json:"founded_year,omitempty"`
}

type VehicleModelUpdate struct {
	Name           *string `json:"name,omitempty"`
	Year           *int    `json:"year,omitempty"`
	VehicleBrandID *string `json:"vehicle_brand_id,omitempty"`
	VehicleTypeID  *int    `json:"vehicle_type_id,omitempty"`
	Feature        *string `json:"feature,omitempty"`
}

type VehicleBasic struct {
	ID             int    `json:"id"`
	VehicleTypeID  int    `json:"vehicle_type_id"`
	VehicleBrandID int    `json:"vehicle_brand_id"`
	Numberplate    string `json:"numberplate"`
}
