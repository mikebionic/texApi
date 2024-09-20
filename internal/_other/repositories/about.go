package repositories

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	db "texApi/database"
	"texApi/internal/_other/queries"
	"texApi/internal/_other/schemas/request"
	"texApi/internal/_other/schemas/response"
)

func GetAboutUsAll() ([]response.AboutUs, error) {
	var aboutUsAll []response.AboutUs

	err := pgxscan.Select(
		context.Background(), db.DB,
		&aboutUsAll, queries.GetAboutUsAll,
	)

	if err != nil {
		return nil, err
	}

	return aboutUsAll, nil
}

func GetAboutUs(aboutUsID int) response.AboutUs {
	var aboutUs response.AboutUs

	db.DB.QueryRow(
		context.Background(), queries.GetAboutUs, aboutUsID,
	).Scan(
		&aboutUs.ID,
		&aboutUs.Text,
		&aboutUs.IsActive,
	)

	return aboutUs
}

func GetAboutUsForUser() response.AboutUs {
	var aboutUs response.AboutUs

	db.DB.QueryRow(context.Background(), queries.GetAboutUsForUser).Scan(
		&aboutUs.ID,
		&aboutUs.Text,
		&aboutUs.IsActive,
	)

	return aboutUs
}

func CreateAboutUs(aboutUs request.CreateAboutUs) (int, error) {
	var aboutUsID int

	db.DB.QueryRow(context.Background(), queries.CreateAboutUs).Scan(&aboutUsID)

	_, err := db.DB.Exec(
		context.Background(), queries.CreateAboutUsTranslates,
		aboutUs.Text.TK, 1,
		aboutUs.Text.RU, 2,
		aboutUs.Text.EN, 3,
		aboutUsID,
	)

	if err != nil {
		return aboutUsID, err
	}

	return aboutUsID, nil
}

func UpdateAboutUs(aboutUs request.UpdateAboutUs) error {
	_, err := db.DB.Exec(
		context.Background(), queries.UpdateAboutUsStatus,
		aboutUs.IsActive, aboutUs.ID,
	)

	if err != nil {
		return err
	}

	_, errTK := db.DB.Exec(
		context.Background(), queries.UpdateAboutUsText,
		aboutUs.Text.TK, 1, aboutUs.ID,
	)

	if errTK != nil {
		return errTK
	}

	_, errRU := db.DB.Exec(
		context.Background(), queries.UpdateAboutUsText,
		aboutUs.Text.RU, 2, aboutUs.ID,
	)

	if errRU != nil {
		return errTK
	}

	_, errEN := db.DB.Exec(
		context.Background(), queries.UpdateAboutUsText,
		aboutUs.Text.EN, 3, aboutUs.ID,
	)

	if errEN != nil {
		return errTK
	}

	return nil
}

func DeleteAboutUs(aboutUsID int) error {
	_, err := db.DB.Exec(context.Background(), queries.DeleteAboutUs, aboutUsID)

	if err != nil {
		return err
	}

	return nil
}
