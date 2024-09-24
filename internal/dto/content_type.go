package dto

import "github.com/google/uuid"

type ContentType struct {
	ID          int               `json:"id"`
	UUID        uuid.UUID         `json:"uuid"`
	Name        string            `json:"name"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	ContentData []ContentResponse `json:"content_data"`
}

type ContentTypeWithContent struct {
	//ID          int               `db:"content_type_id" json:"id"`
	UUID        uuid.UUID         `db:"content_type_uuid" json:"uuid"`
	Name        string            `db:"content_type_name" json:"name"`
	Title       string            `db:"content_type_title" json:"title"`
	Description string            `db:"content_type_description" json:"description"`
	Contents    []ContentResponse `json:"contents"`
}
