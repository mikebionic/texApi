package dto

import "time"

type MediaCreate struct {
	ID        int    `json:"id" db:"id"`
	UUID      string `json:"uuid" db:"uuid"`
	UserID    int    `json:"user_id"`
	CompanyID int    `json:"company_id,omitempty"`

	// Content metadata
	MediaType   string  `json:"media_type"`
	Context     string  `json:"context"`
	ContextID   *int    `json:"context_id,omitempty"`
	ContextUUID *string `json:"context_uuid,omitempty"`

	// File information
	Filename   string  `json:"filename"`
	FilePath   *string `json:"file_path,omitempty"`
	ThumbPath  *string `json:"thumb_path,omitempty"`
	ThumbFn    *string `json:"thumb_fn,omitempty"`
	OriginalFn string  `json:"original_fn"`

	// Media metadata
	MimeType *string `json:"mime_type,omitempty"`
	FileSize *int64  `json:"file_size,omitempty"`
	Duration *int    `json:"duration,omitempty"`
	Width    *int    `json:"width,omitempty"`
	Height   *int    `json:"height,omitempty"`

	// Additional metadata fields
	Meta  *string `json:"meta,omitempty"`
	Meta2 *string `json:"meta2,omitempty"`
	Meta3 *string `json:"meta3,omitempty"`
}

type MediaMain struct {
	MediaCreate

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Active    int       `json:"active" db:"active"`
	Deleted   int       `json:"deleted" db:"deleted"`

	URL      string `json:"url,omitempty"`
	ThumbURL string `json:"thumb_url,omitempty"`
}
