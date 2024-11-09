package dto

type Driver struct {
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
}

type DriverDetails struct {
	ID               int            `json:"id"`
	UUID             string         `json:"uuid"`
	CompanyID        int            `json:"company_id"`
	FirstName        string         `json:"first_name"`
	LastName         string         `json:"last_name"`
	PatronymicName   string         `json:"patronymic_name"`
	Phone            string         `json:"phone"`
	Email            string         `json:"email"`
	Featured         int            `json:"featured"`
	Rating           int            `json:"rating"`
	Partner          int            `json:"partner"`
	SuccessfulOps    int            `json:"successful_ops"`
	ImageURL         string         `json:"image_url"`
	CreatedAt        string         `json:"created_at"`
	UpdatedAt        string         `json:"updated_at"`
	Active           int            `json:"active"`
	Company          *CompanyBasic  `json:"company,omitempty"`
	AssignedVehicles []VehicleBasic `json:"assigned_vehicles,omitempty"`
}

type DriverCreate struct {
	CompanyID      int    `json:"company_id" binding:"required"`
	FirstName      string `json:"first_name" binding:"required"`
	LastName       string `json:"last_name" binding:"required"`
	PatronymicName string `json:"patronymic_name"`
	Phone          string `json:"phone" binding:"required"`
	Email          string `json:"email"`
	ImageURL       string `json:"image_url"`
}

type DriverUpdate struct {
	FirstName      *string `json:"first_name,omitempty"`
	LastName       *string `json:"last_name,omitempty"`
	PatronymicName *string `json:"patronymic_name,omitempty"`
	Phone          *string `json:"phone,omitempty"`
	Email          *string `json:"email,omitempty"`
	Featured       *int    `json:"featured,omitempty"`
	Rating         *int    `json:"rating,omitempty"`
	Partner        *int    `json:"partner,omitempty"`
	ImageURL       *string `json:"image_url,omitempty"`
	Active         *int    `json:"active,omitempty"`
}
