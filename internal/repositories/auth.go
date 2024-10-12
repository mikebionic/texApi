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

func ManageToken(id int, token, action string) error {
	switch action {
	case "create":
		_, err := db.DB.Exec(
			context.Background(), `UPDATE tbl_user SET refresh_token = $1 WHERE id = $2`,
			token,
			id,
		)
		return err

	case "validate":
		var currentToken string
		err := db.DB.QueryRow(context.Background(), `SELECT refresh_token FROM tbl_user WHERE id = $1`, id).Scan(&currentToken)
		if err != nil {
			return err
		}

		if currentToken == token {
			return nil
		} else {
			return fmt.Errorf("token is invalid")
		}

	default:
		return fmt.Errorf("invalid actionб must be create or validate")
	}
}

func CreateUser(user dto.CreateUser) (int, error) {
	var id int

	err := db.DB.QueryRow(
		context.Background(), queries.CreateUser,
		user.Username,
		user.Password,
		user.Email,
		user.FirstName,
		user.LastName,
		user.NickName,
		user.AvatarURL,
		user.Phone,
		user.InfoPhone,
		user.Address,
		user.RoleID,
		user.SubroleID,
		user.Verified,
		user.Active,
		user.OauthProvider,
		user.OauthUserID,
		user.OauthLocation,
		user.OauthAccessToken,
		user.OauthAccessTokenSecret,
		user.OauthRefreshToken,
		user.OauthIDToken,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
