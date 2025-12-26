package dto

import "time"

type DriverShort struct {
	ID             int    `json:"id"`
	UUID           string `json:"uuid"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	PatronymicName string `json:"patronymic_name"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
	Featured       int    `json:"featured"`
	Rating         int    `json:"rating"`
	Partner        int    `json:"partner"`
	SuccessfulOps  int    `json:"successful_ops"`
	ImageURL       string `json:"image_url"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type DriverDetails struct {
	DriverCreate
	UserID           *int             `json:"user_id"`
	Company          *CompanyBasic    `json:"company,omitempty"`
	AssignedVehicles []VehicleBasic   `json:"assigned_vehicles,omitempty"`
	UserCredentials  *UserCredentials `json:"user,omitempty"`
}

type UserCredentials struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type DriverCreate struct {
	ID             int       `json:"id"`
	UUID           string    `json:"uuid"`
	CompanyID      int       `json:"company_id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	PatronymicName string    `json:"patronymic_name"`
	Phone          string    `json:"phone"`
	Email          string    `json:"email"`
	Featured       int       `json:"featured"`
	Rating         int       `json:"rating"`
	Partner        int       `json:"partner"`
	SuccessfulOps  int       `json:"successful_ops"`
	ImageURL       string    `json:"image_url"`
	ViewCount      int       `json:"view_count"`
	Meta           string    `json:"meta"`
	Meta2          string    `json:"meta2"`
	Meta3          string    `json:"meta3"`
	Available      int       `json:"available"`
	BlockReason    *string   `json:"block_reason"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Active         int       `json:"active"`
	Deleted        int       `json:"deleted"`
	TotalCount     string    `json:"total_count"`
}

type DriverUpdate struct {
	CompanyID      *int    `json:"company_id,omitempty"`
	FirstName      *string `json:"first_name,omitempty"`
	LastName       *string `json:"last_name,omitempty"`
	PatronymicName *string `json:"patronymic_name,omitempty"`
	Phone          *string `json:"phone,omitempty"`
	Email          *string `json:"email,omitempty"`
	ImageURL       *string `json:"image_url,omitempty"`
	Meta           *string `json:"meta,omitempty"`
	Meta2          *string `json:"meta2,omitempty"`
	Meta3          *string `json:"meta3,omitempty"`
	BlockReason    *string `json:"block_reason"`
	Active         *int    `json:"active,omitempty"`
	Deleted        *int    `json:"deleted,omitempty"`
}
