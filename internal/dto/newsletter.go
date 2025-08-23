package dto

import "time"

type NewsletterMain struct {
	ID          int    `json:"id"`
	UUID        string `json:"uuid"`
	Email       string `json:"email"`
	Status      string `json:"status"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Frequency   string `json:"frequency"`
	IPAddress   string `json:"ip_address"`
	UserAgent   string `json:"user_agent"`
	ReferrerURL string `json:"referrer_url"`
	Meta        string `json:"meta"`
	Meta2       string `json:"meta2"`
	Meta3       string `json:"meta3"`
	Active      int    `json:"active"`
	Deleted     int    `json:"deleted"`
}

type Newsletter struct {
	NewsletterMain
	SubscribedAt   *time.Time `json:"subscribed_at"`
	UnsubscribedAt *time.Time `json:"unsubscribed_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	TotalCount     int        `json:"total_count"`
}

type NewsletterUpdate struct {
	Email       *string `json:"email,omitempty"`
	Status      *string `json:"status,omitempty"`
	FirstName   *string `json:"first_name,omitempty"`
	LastName    *string `json:"last_name,omitempty"`
	Frequency   *string `json:"frequency,omitempty"`
	IPAddress   *string `json:"ip_address,omitempty"`
	UserAgent   *string `json:"user_agent,omitempty"`
	ReferrerURL *string `json:"referrer_url,omitempty"`
	Meta        *string `json:"meta,omitempty"`
	Meta2       *string `json:"meta2,omitempty"`
	Meta3       *string `json:"meta3,omitempty"`
	Active      *int    `json:"active,omitempty"`
	Deleted     *int    `json:"deleted,omitempty"`
}

type NewsletterSubscribe struct {
	Email       string `json:"email" binding:"required,email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Frequency   string `json:"frequency"`
	IPAddress   string `json:"ip_address"`
	UserAgent   string `json:"user_agent"`
	ReferrerURL string `json:"referrer_url"`
	Meta        string `json:"meta"`
	Meta2       string `json:"meta2"`
	Meta3       string `json:"meta3"`
}
