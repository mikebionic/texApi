package dto

type DriverCreate struct {
	ID             int    `json:"id"`
	CompanyID      int    `json:"company_id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	PatronymicName string `json:"patronymic_name"`
	Phone          string `json:"phone"`
	Email          string `json:"email"`
	AvatarURL      string `json:"avatar_url"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	Active         int    `json:"active"`
	Deleted        int    `json:"deleted"`
}

type DriverUpdate struct {
	FirstName      *string `json:"first_name,omitempty"`
	LastName       *string `json:"last_name,omitempty"`
	PatronymicName *string `json:"patronymic_name,omitempty"`
	Phone          *string `json:"phone,omitempty"`
	Email          *string `json:"email,omitempty"`
	AvatarURL      *string `json:"avatar_url,omitempty"`
}
