package repositories

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	db "texApi/database"
	"texApi/internal/dto"
	"texApi/internal/queries"
)

func GetContentTypes(withContent int) ([]dto.ContentType, error) {
	var contentTypes []dto.ContentType
	if withContent > 0 {
		fmt.Println("Withcontent")
		var results []struct {
			ContentType dto.ContentType
			Content     dto.ContentResponse
		}
		err := pgxscan.Select(context.Background(), db.DB, &results, queries.GetContentTypesWithContent)
		if err != nil {
			return nil, err
		}

		contentTypesMap := make(map[int]*dto.ContentType)
		for _, row := range results {
			if _, exists := contentTypesMap[row.ContentType.ID]; !exists {
				contentTypesMap[row.ContentType.ID] = &row.ContentType
				contentTypesMap[row.ContentType.ID].ContentData = []dto.ContentResponse{}
			}

			if row.Content.ID != 0 { // Ensure we only append non-null content rows
				contentTypesMap[row.ContentType.ID].ContentData = append(contentTypesMap[row.ContentType.ID].ContentData, row.Content)
			}
		}

		for _, ct := range contentTypesMap {
			contentTypes = append(contentTypes, *ct)
		}
	} else {
		err := pgxscan.Select(
			context.Background(), db.DB,
			&contentTypes, queries.GetContentTypes)
		if err != nil {
			return nil, err
		}
	}
	return contentTypes, nil
}
