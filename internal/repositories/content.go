package repositories

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	db "texApi/database"
	"texApi/internal/dto"
	"texApi/internal/queries"
)

func GetContents(ctID int) ([]dto.ContentResponse, error) {
	var contents []dto.ContentResponse
	var err error
	stmt := queries.GetContents
	if ctID > 0 {
		stmt += ` AND c.content_type_id=$1`
		err = pgxscan.Select(
			context.Background(), db.DB,
			&contents, stmt, ctID,
		)
	} else {
		err = pgxscan.Select(
			context.Background(), db.DB,
			&contents, stmt,
		)
	}
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func GetContent(id int) dto.ContentResponse {
	var content dto.ContentResponse

	db.DB.QueryRow(
		context.Background(), queries.GetContent, id,
	).Scan(
		&content.ID,
		&content.UUID,
		&content.LangID,
		&content.ContentTypeID,
		&content.Title,
		&content.Slogan,
		&content.Subtitle,
		&content.Description,
		&content.Count,
		&content.CountType,
		&content.ImageURL,
		&content.VideoURL,
		&content.Step,
		&content.CreatedAt,
		&content.UpdatedAt,
		&content.Active,
		&content.Deleted,
	)

	return content
}

func CreateContent(content dto.CreateContent) int {
	var id int

	db.DB.QueryRow(
		context.Background(), queries.CreateContent,
		content.LangID,
		content.ContentTypeID,
		content.Title,
		content.Slogan,
		content.Subtitle,
		content.Description,
		content.Count,
		content.CountType,
		content.ImageURL,
		content.VideoURL,
		content.Step,
		content.Active,
	).Scan(&id)

	return id
}

func UpdateContent(content dto.CreateContent, id int) (updatedID int, err error) {
	err = db.DB.QueryRow(
		context.Background(), queries.UpdateContent,
		content.LangID,
		content.ContentTypeID,
		content.Title,
		content.Slogan,
		content.Subtitle,
		content.Description,
		content.Count,
		content.CountType,
		content.ImageURL,
		content.VideoURL,
		content.Step,
		content.Active,
		id,
	).Scan(&updatedID)
	
	return
}

func DeleteContent(id int) error {
	_, err := db.DB.Exec(context.Background(), queries.DeleteContent, id)
	if err != nil {
		return err
	}

	return nil
}
