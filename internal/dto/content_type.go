package dto

import "github.com/google/uuid"

type ContentType struct {
	ID          int               `json:"id"`
	UUID        uuid.UUID         `json:"uuid"`
	Name        string            `json:"name"`
	Title       string            `json:"title"`
	TitleRu     string            `json:"title_ru"`
	Description string            `json:"description"`
	ParentID    int               `json:"parent_id"`
	ParentName  string            `json:"parent_name"`
	ContentData []ContentResponse `json:"content_data"`
}

type ContentTypeWithContent struct {
	ID          int               `db:"content_type_id" json:"id"`
	UUID        uuid.UUID         `db:"content_type_uuid" json:"uuid"`
	Name        string            `db:"content_type_name" json:"name"`
	Title       string            `db:"content_type_title" json:"title"`
	TitleRu     string            `db:"content_type_title_ru" json:"title_ru"`
	Description string            `db:"content_type_description" json:"description"`
	ParentID    string            `db:"content_type_parent_id" json:"parent_id"`
	ParentName  string            `db:"content_type_parent_name" json:"parent_name"`
	Contents    []ContentResponse `json:"contents"`
}
