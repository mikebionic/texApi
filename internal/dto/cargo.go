package dto

import "time"

type CargoMain struct {
	ID              int    `json:"id"`
	UUID            string `json:"uuid"`
	CompanyID       int    `json:"company_id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Info            string `json:"info"`
	Qty             int    `json:"qty"`
	Weight          int    `json:"weight"`
	WeightType      string `json:"weight_type"`
	Meta            string `json:"meta"`
	Meta2           string `json:"meta2"`
	Meta3           string `json:"meta3"`
	VehicleTypeID   int    `json:"vehicle_type_id"`
	PackagingTypeID int    `json:"packaging_type_id"`
	GPS             int    `json:"gps"`
	Photo1URL       string `json:"photo1_url"`
	Photo2URL       string `json:"photo2_url"`
	Photo3URL       string `json:"photo3_url"`
	Docs1URL        string `json:"docs1_url"`
	Docs2URL        string `json:"docs2_url"`
	Docs3URL        string `json:"docs3_url"`
	Note            string `json:"note"`
	Active          int    `json:"active"`
	Deleted         int    `json:"deleted"`
}
type Cargo struct {
	CargoMain
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	TotalCount int       `json:"total_count"`
}

type CargoUpdate struct {
	CompanyID       *int    `json:"company_id,omitempty"`
	Name            *string `json:"name,omitempty"`
	Description     *string `json:"description,omitempty"`
	Info            *string `json:"info,omitempty"`
	Qty             *int    `json:"qty,omitempty"`
	Weight          *int    `json:"weight,omitempty"`
	WeightType      *string `json:"weight_type"`
	Meta            *string `json:"meta,omitempty"`
	Meta2           *string `json:"meta2,omitempty"`
	Meta3           *string `json:"meta3,omitempty"`
	VehicleTypeID   *int    `json:"vehicle_type_id,omitempty"`
	PackagingTypeID *int    `json:"packaging_type_id,omitempty"`
	GPS             *int    `json:"gps,omitempty"`
	Photo1URL       *string `json:"photo1_url,omitempty"`
	Photo2URL       *string `json:"photo2_url,omitempty"`
	Photo3URL       *string `json:"photo3_url,omitempty"`
	Docs1URL        *string `json:"docs1_url,omitempty"`
	Docs2URL        *string `json:"docs2_url,omitempty"`
	Docs3URL        *string `json:"docs3_url,omitempty"`
	Note            *string `json:"note,omitempty"`
	Active          *int    `json:"active,omitempty"`
	Deleted         *int    `json:"deleted,omitempty"`
}

type CargoDetailed struct {
	Cargo
	Company       *CompanyCreate         `json:"company"`
	VehicleType   *VehicleType           `json:"vehicle_type"`
	PackagingType *PackagingTypeResponse `json:"packaging_type"`
}
