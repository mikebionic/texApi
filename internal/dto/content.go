package dto

import "github.com/google/uuid"

type ContentResponse struct {
	ID            int       `json:"id"`
	UUID          uuid.UUID `json:"uuid"`
	LangID        int       `json:"lang_id"`
	ContentTypeID int       `json:"content_type_id"`
	Title         string    `json:"title"`
	Subtitle      string    `json:"subtitle"`
	Description   string    `json:"description"`
	ImageURL      string    `json:"image_url"`
	VideoURL      string    `json:"video_url"`
	Step          int       `json:"step"`
	CreatedAt     string    `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
	Deleted       int       `json:"deleted"`
}

type CreateContent struct {
	LangID        int    `json:"lang_id" binding:"gt=0"`
	ContentTypeID int    `json:"content_type_id" binding:"gt=0"`
	Title         string `json:"title"`
	Subtitle      string `json:"subtitle"`
	Description   string `json:"description"`
	ImageURL      string `json:"image_url"`
	VideoURL      string `json:"video_url"`
	Step          int    `json:"step"`
}
