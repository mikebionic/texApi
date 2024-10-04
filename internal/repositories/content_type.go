package repositories

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	db "texApi/database"
	"texApi/internal/dto"
	"texApi/internal/queries"
)

func GetContentTypes(withContent, langID int) ([]dto.ContentTypeWithContents, error) {
	var contentTypes []dto.ContentTypeWithContents
	var stmt string
	var err error

	switch withContent {
	case 1:
		stmt = queries.GetContentTypesWithContent
		err = pgxscan.Select(context.Background(), db.DB, &contentTypes, stmt, langID)
	default:
		stmt = queries.GetContentTypes
		err = pgxscan.Select(context.Background(), db.DB, &contentTypes, stmt)
	}

	if err != nil {
		return nil, err
	}
	return contentTypes, nil
}
