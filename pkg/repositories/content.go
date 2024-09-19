package repositories

import (
	"context"
	db "texApi/database"
	"texApi/pkg/queries"
	"texApi/pkg/schemas/response"

	"github.com/georgysavva/scany/v2/pgxscan"
)

func GetContents() ([]response.Content, error) {
	var contents []response.Content
	err := pgxscan.Select(
		context.Background(), db.DB,
		&contents, queries.GetContents,
	)
	if err != nil {
		return nil, err
	}
	return contents, nil
}
