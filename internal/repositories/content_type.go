package repositories

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	db "texApi/database"
	"texApi/internal/dto"
	"texApi/internal/queries"
)

func GetContentTypes(withContent int) ([]dto.ContentType, error) {
	var contentTypes []dto.ContentType
	if withContent > 0 {
		var results []struct {
			ContentType dto.ContentTypeWithContent
			Content     dto.ContentResponse
		}
		err := pgxscan.Select(context.Background(), db.DB, &results, queries.GetContentTypesWithContent)
		if err != nil {
			return nil, err
		}

		//contentTypesMap := make(map[int]*dto.ContentTypeWithContent)
		//for _, row := range results {
		//	fmt.Println(row)
		//	if _, exists := contentTypesMap[row.ContentType.ID]; !exists {
		//		contentTypesMap[row.ContentType.ID] = &row.ContentType
		//		contentTypesMap[row.ContentType.ID].Contents = []dto.ContentResponse{}
		//	}
		//
		//	if row.Content.ID != 0 {
		//		contentTypesMap[row.ContentType.ID].Contents = append(
		//			contentTypesMap[row.ContentType.ID].Contents,
		//			row.Content,
		//		)
		//	}
		//}

		//for _, ct := range contentTypesMap {
		//	contentTypes = append(contentTypes, *ct)
		//}
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
