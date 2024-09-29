package dto

import "github.com/google/uuid"

type ContentResponse struct {
	ID            int       `db:"id" json:"id"`
	UUID          uuid.UUID `db:"uuid" json:"uuid"`
	LangID        int       `db:"lang_id" json:"lang_id"`
	ContentTypeID int       `db:"content_type_id" json:"content_type_id"`
	Title         string    `db:"title" json:"title"`
	Slogan        string    `db:"slogan" json:"slogan"`
	Subtitle      string    `db:"subtitle" json:"subtitle"`
	Description   string    `db:"description" json:"description"`
	ImageURL      string    `db:"image_url" json:"image_url"`
	VideoURL      string    `db:"video_url" json:"video_url"`
	Step          int       `db:"step" json:"step"`
	CreatedAt     string    `db:"created_at" json:"created_at"`
	UpdatedAt     string    `db:"updated_at" json:"updated_at"`
	Active        int       `db:"active" json:"active"`
	Deleted       int       `db:"deleted" json:"deleted"`
}

type CreateContent struct {
	LangID        int    `json:"lang_id" binding:"gt=0"`
	ContentTypeID int    `json:"content_type_id" binding:"gt=0"`
	Title         string `json:"title"`
	Slogan        string `json:"slogan"`
	Subtitle      string `json:"subtitle"`
	Description   string `json:"description"`
	ImageURL      string `json:"image_url"`
	VideoURL      string `json:"video_url"`
	Step          int    `json:"step"`
	Active        int    `json:"active"`
}
