package dto

import "github.com/google/uuid"

type ContentResponse struct {
	ID            int       `json:"id"`
	UUID          uuid.UUID `json:"uuid"`
	LangID        int       `json:"lang_id"`
	ContentTypeID int       `json:"content_type_id"`
	Title         string    `json:"title"`
	Slogan        string    `json:"slogan"`
	Subtitle      string    `json:"subtitle"`
	Description   string    `json:"description"`
	Count         int       `json:"count"`
	CountType     string    `json:"count_type"`
	ImageURL      string    `json:"image_url"`
	VideoURL      string    `json:"video_url"`
	Step          int       `json:"step"`
	CreatedAt     string    `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
	Active        int       `json:"active"`
	Deleted       int       `json:"deleted"`
}

type CreateContent struct {
	LangID        *int    `json:"lang_id,omitempty"`
	ContentTypeID *int    `json:"content_type_id,omitempty"`
	Title         *string `json:"title,omitempty"`
	Slogan        *string `json:"slogan,omitempty"`
	Subtitle      *string `json:"subtitle,omitempty"`
	Description   *string `json:"description,omitempty"`
	Count         *int    `json:"count,omitempty"`
	CountType     *string `json:"count_type,omitempty"`
	ImageURL      *string `json:"image_url,omitempty"`
	VideoURL      *string `json:"video_url,omitempty"`
	Step          *int    `json:"step,omitempty"`
	Active        *int    `json:"active,omitempty"`
}
