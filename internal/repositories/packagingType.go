package repositories

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	db "texApi/database"
	"texApi/internal/dto"
	"texApi/internal/queries"
)

func GetPackagingTypes() ([]dto.PackagingTypeResponse, error) {
	var types []dto.PackagingTypeResponse
	err := pgxscan.Select(
		context.Background(),
		db.DB,
		&types,
		queries.GetPackagingTypes,
	)
	if err != nil {
		return nil, err
	}
	return types, nil
}

func GetPackagingType(id int) (dto.PackagingTypeResponse, error) {
	var packagingType dto.PackagingTypeResponse
	err := pgxscan.Get(
		context.Background(),
		db.DB,
		&packagingType,
		queries.GetPackagingType,
		id,
	)
	if err != nil {
		return dto.PackagingTypeResponse{}, err
	}
	return packagingType, nil
}

func CreatePackagingType(pt dto.CreatePackagingType) (int, error) {
	var id int
	err := db.DB.QueryRow(
		context.Background(),
		queries.CreatePackagingType,
		pt.NameRu, pt.NameEn, pt.NameTk,
		pt.CategoryRu, pt.CategoryEn, pt.CategoryTk,
		pt.Material, pt.Dimensions, pt.Weight,
		pt.DescriptionRu, pt.DescriptionEn, pt.DescriptionTk,
		pt.Active,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func UpdatePackagingType(pt dto.CreatePackagingType, id int) (int, error) {
	var updatedID int
	err := db.DB.QueryRow(
		context.Background(),
		queries.UpdatePackagingType,
		pt.NameRu, pt.NameEn, pt.NameTk,
		pt.CategoryRu, pt.CategoryEn, pt.CategoryTk,
		pt.Material, pt.Dimensions, pt.Weight,
		pt.DescriptionRu, pt.DescriptionEn, pt.DescriptionTk,
		pt.Active,
		id,
	).Scan(&updatedID)
	return updatedID, err
}

func DeletePackagingType(id int) error {
	_, err := db.DB.Exec(context.Background(), queries.DeletePackagingType, id)
	return err
}
