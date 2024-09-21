package dto

import "github.com/google/uuid"

type ContentType struct {
	ID          int       `json:"id"`
	UUID        uuid.UUID `json:"uuid"`
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}
