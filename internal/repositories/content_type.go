package repositories

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	db "texApi/database"
	"texApi/internal/dto"
	"texApi/internal/queries"
)

func GetContentTypes() ([]dto.ContentType, error) {
	var contentTypes []dto.ContentType
	err := pgxscan.Select(
		context.Background(), db.DB,
		&contentTypes, queries.GetContentTypes)
	if err != nil {
		return nil, err
	}
	return contentTypes, nil
}
