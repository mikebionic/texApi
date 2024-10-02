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
	Contents    []ContentResponse `json:"content_data"`
}

type ContentTypeWithContents struct {
	ID          int               `db:"id" json:"id"`
	UUID        uuid.UUID         `db:"uuid" json:"uuid"`
	Name        string            `db:"name" json:"name"`
	Title       string            `db:"title" json:"title"`
	TitleRu     string            `db:"title_ru" json:"title_ru"`
	Description string            `db:"description" json:"description"`
	ParentID    int               `db:"parent_id" json:"parent_id"`
	ParentName  string            `db:"parent_name" json:"parent_name"`
	Contents    []ContentResponse `db:"contents" json:"contents"`
}
