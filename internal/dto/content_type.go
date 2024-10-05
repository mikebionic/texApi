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
	ID          int               `json:"id"`
	UUID        uuid.UUID         `json:"uuid"`
	Name        string            `json:"name"`
	Title       string            `json:"title"`
	TitleRu     string            `json:"title_ru"`
	Description string            `json:"description"`
	ParentID    int               `json:"parent_id"`
	ParentName  string            `json:"parent_name"`
	Contents    []ContentResponse `json:"contents"`
}
