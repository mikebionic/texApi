package repositories

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	db "texApi/database"
	"texApi/internal/dto"
	"texApi/internal/queries"
)

func GetContentTypes(withContent, langID, ctID int) ([]dto.ContentTypeWithContents, error) {
	var contentTypes []dto.ContentTypeWithContents
	var stmt string
	var err error

	switch withContent {
	case 1:
		stmt = queries.GetContentTypesWithContent
		err = pgxscan.Select(context.Background(), db.DB, &contentTypes, stmt, langID, ctID)
	default:
		stmt = queries.GetContentTypes
		err = pgxscan.Select(context.Background(), db.DB, &contentTypes, stmt)
	}

	if err != nil {
		return nil, err
	}
	if len(contentTypes) == 0 {
		return nil, fmt.Errorf("not found, empty slice")
	}
	return contentTypes, nil
}
