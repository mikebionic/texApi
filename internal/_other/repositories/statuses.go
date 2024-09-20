package repositories

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	db "texApi/database"
	"texApi/internal/_other/queries"
	"texApi/internal/_other/schemas/response"
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
