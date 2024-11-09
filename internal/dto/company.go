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
