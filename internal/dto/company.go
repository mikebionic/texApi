package dto

import "time"

type CompanyDetails struct {
	CompanyCreate
	Drivers  []Driver  `json:"drivers,omitempty"`
	Vehicles []Vehicle `json:"vehicles,omitempty"`
}

type CompanyCreate struct {
	ID             int       `json:"id"`
	UUID           string    `json:"uuid"`
	UserID         int       `json:"user_id"`
	RoleID         int       `json:"role_id"`
	CompanyName    string    `json:"company_name"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	PatronymicName string    `json:"patronymic_name"`
	Phone          string    `json:"phone"`
	Phone2         string    `json:"phone2"`
	Phone3         string    `json:"phone3"`
	Email          string    `json:"email"`
	Email2         string    `json:"email2"`
	Email3         string    `json:"email3"`
	Meta           string    `json:"meta"`
	Meta2          string    `json:"meta2"`
	Meta3          string    `json:"meta3"`
	Address        string    `json:"address"`
	Country        string    `json:"country"`
	CountryID      int       `json:"country_id"`
	CityID         int       `json:"city_id"`
	ImageURL       string    `json:"image_url"`
	Entity         string    `json:"entity"`
	Featured       int       `json:"featured"`
	Rating         int       `json:"rating"`
	Partner        int       `json:"partner"`
	SuccessfulOps  int       `json:"successful_ops"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Active         int       `json:"active"`
	Deleted        int       `json:"deleted"`
}

type CompanyUpdate struct {
	UserID         *int    `json:"user_id,omitempty"`
	RoleID         *int    `json:"role_id,omitempty"`
	CompanyName    *string `json:"company_name,omitempty"`
	FirstName      *string `json:"first_name,omitempty"`
	LastName       *string `json:"last_name,omitempty"`
	PatronymicName *string `json:"patronymic_name,omitempty"`
	Phone          *string `json:"phone,omitempty"`
	Phone2         *string `json:"phone2,omitempty"`
	Phone3         *string `json:"phone3,omitempty"`
	Email          *string `json:"email,omitempty"`
	Email2         *string `json:"email2,omitempty"`
	Email3         *string `json:"email3,omitempty"`
	Meta           *string `json:"meta,omitempty"`
	Meta2          *string `json:"meta2,omitempty"`
	Meta3          *string `json:"meta3,omitempty"`
	Address        *string `json:"address,omitempty"`
	Country        *string `json:"country,omitempty"`
	CountryID      *int    `json:"country_id,omitempty"`
	CityID         *int    `json:"city_id,omitempty"`
	ImageURL       *string `json:"image_url,omitempty"`
	Entity         *string `json:"entity,omitempty"`
	Featured       *int    `json:"featured,omitempty"`
	Rating         *int    `json:"rating,omitempty"`
	Partner        *int    `json:"partner,omitempty"`
	SuccessfulOps  *int    `json:"successful_ops,omitempty"`
	Active         *int    `json:"active,omitempty"`
	Deleted        *int    `json:"deleted,omitempty"`
}

// Basic DTOs for nested responses
type CompanyBasic struct {
	ID          int    `json:"id"`
	CompanyName string `json:"company_name"`
	Country     string `json:"country"`
}
