package dto

//type CompanyGet struct {
//	ID      int    `json:"id"`
//	UserID  int    `json:"user_id"`
//	Name    string `json:"name"`
//	Address string `json:"address"`
//	Phone   string `json:"phone"`
//	Email   string `json:"email"`
//	LogoURL string `json:"logo_url"`
//}

type CompanyDetails struct {
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
	CreatedAt      string    `json:"created_at"`
	UpdatedAt      string    `json:"updated_at"`
	Active         int       `json:"active"`
	Drivers        []Driver  `json:"drivers,omitempty"`
	Vehicles       []Vehicle `json:"vehicles,omitempty"`
}

type PaginatedResponse struct {
	Total   int         `json:"total"`
	Page    int         `json:"page"`
	PerPage int         `json:"per_page"`
	Data    interface{} `json:"data"`
}

type CompanyCreate struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	Country   string `json:"country"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	LogoURL   string `json:"logo_url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Active    int    `json:"active"`
	Deleted   int    `json:"deleted"`
}

type CompanyUpdate struct {
	Name    *string `json:"name,omitempty"`
	Address *string `json:"address,omitempty"`
	Country *string `json:"country,omitempty"`
	Phone   *string `json:"phone,omitempty"`
	Email   *string `json:"email,omitempty"`
	LogoURL *string `json:"logo_url,omitempty"`
}

// Basic DTOs for nested responses
type CompanyBasic struct {
	ID          int    `json:"id"`
	CompanyName string `json:"company_name"`
	Country     string `json:"country"`
}
