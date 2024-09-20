package repositories

import (
	"context"
	db "texApi/database"
	"texApi/internal/_other/queries"
	"texApi/internal/_other/schemas/response"
)

func GetAdmin(phone string) response.Login {
	var admin response.Login

	db.DB.QueryRow(
		context.Background(), queries.GetAdmin, phone,
	).Scan(&admin.ID, &admin.Phone, &admin.Password)

	return admin
}

func GetUserForLogin(phone string) response.Login {
	var user response.Login

	db.DB.QueryRow(
		context.Background(), queries.GetUserForLogin, phone,
	).Scan(
		&user.ID,
		&user.Fullname,
		&user.Phone,
		&user.Address,
		&user.Password,
		&user.Subscription,
	)

	return user
}

func GetUserMe(userID int) response.Login {
	var user response.Login

	db.DB.QueryRow(
		context.Background(), queries.GetUserMe, userID,
	).Scan(
		&user.ID,
		&user.Fullname,
		&user.Phone,
		&user.Address,
		&user.Subscription,
	)

	return user
}

func GetWorkerForLogin(phone string) response.Login {
	var worker response.Login

	db.DB.QueryRow(
		context.Background(), queries.GetWorkerForLogin, phone,
	).Scan(&worker.ID, &worker.Phone, &worker.Password)

	return worker
}
