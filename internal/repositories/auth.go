package repositories

import (
	"context"
	db "texApi/database"
	"texApi/internal/_other/schemas/response"
)

func GetUserByPhone(phone string) response.Login {
	var user response.Login
	db.DB.QueryRow(
		context.Background(), "SELECT id, fullname, phone, address, password FROM tbl_user WHERE phone = $1", phone,
	).Scan(&user.ID, &user.Fullname, &user.Phone, &user.Address, &user.Password)
	return user
}

func GetUserById(userID int) response.Login {
	var user response.Login
	db.DB.QueryRow(
		context.Background(), "SELECT id, fullname, phone, address FROM tbl_user WHERE id = $1", userID,
	).Scan(&user.ID, &user.Fullname, &user.Phone, &user.Address)
	return user
}
