package repositories

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	db "texApi/database"
	"texApi/internal/dto"
	"texApi/internal/queries"
)

func GetUser(username, loginMethod string) (dto.User, error) {
	stmt := queries.GetUser
	switch loginMethod {
	case "phone":
		stmt = stmt + " WHERE phone = $1"
	case "username":
		stmt = stmt + " WHERE username = $1"
	case "email":
		stmt = stmt + " WHERE email = $1"
	}

	var user dto.User
	err := pgxscan.Get(
		context.Background(),
		db.DB,
		&user,
		stmt,
		username,
	)
	if err != nil {
		return user, err
		//fmt.Errorf("error fetching user")
	}

	if user.ID == 0 {
		return user, fmt.Errorf("login failed")
	}

	return user, nil
}

func GetUserById(userID int) dto.User {
	var user dto.User
	_ = pgxscan.Get(
		context.Background(),
		db.DB,
		&user,
		fmt.Sprintf("%s WHERE id = $1", queries.GetUser),
		userID,
	)
	return user
}

func ManageToken(id int, token, action string) (string, error) {

	switch action {
	case "create":
		//_, err := db.DB.Exec(`UPDATE tbl_user SET refresh_token = $1 WHERE id = $2`, token, id)
		//if err != nil {
		//	return "", err
		//}
		return "Token updated successfully", nil

	//case "refresh":
	//	var currentToken string
	//	err := db.DB.QueryRow(`SELECT refresh_token FROM tbl_user WHERE id = $1`, id).Scan(&currentToken)
	//	if err != nil {
	//		return "", err
	//	}
	//
	//	if currentToken == token {
	//		return "Token is valid", nil
	//	} else {
	//		return "", fmt.Errorf("token is invalid")
	//	}

	default:
		return "", fmt.Errorf("invalid action")
	}
}
