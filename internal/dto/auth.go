package dto

import "github.com/google/uuid"

type User struct {
	ID        int       `json:"id"`
	UUID      uuid.UUID `json:"uuid"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	Email     string    `json:"email"`
	Fullname  string    `json:"fullname"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	RoleID    int       `json:"role_id"`
	Verified  int       `json:"verified"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Active    int       `json:"active"`
	Deleted   int       `json:"deleted"`
}
