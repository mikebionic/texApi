package repositories

import (
	"context"
	db "texApi/database"
	"texApi/pkg/queries"
	"texApi/pkg/schemas/response"

	"github.com/georgysavva/scany/v2/pgxscan"
)

func GetStatuses() ([]response.Translate, error) {
	var statuses []response.Translate

	err := pgxscan.Select(
		context.Background(), db.DB,
		&statuses, queries.GetStatuses,
	)

	if err != nil {
		return nil, err
	}

	return statuses, nil
}
