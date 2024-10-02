package repositories

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	db "texApi/database"
	"texApi/internal/dto"
	"texApi/internal/queries"
)

func GetContentTypes(withContent int) ([]dto.ContentTypeWithContents, error) {
	var contentTypes []dto.ContentTypeWithContents

	if withContent > 0 {
		err := pgxscan.Select(context.Background(), db.DB, &contentTypes, queries.GetContentTypesWithContent)
		if err != nil {
			return nil, err
		}
	} else {
		err := pgxscan.Select(context.Background(), db.DB, &contentTypes, queries.GetContentTypes)
		if err != nil {
			return nil, err
		}
	}

	return contentTypes, nil
}
