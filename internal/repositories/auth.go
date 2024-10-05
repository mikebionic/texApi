package repositories

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"golang.org/x/crypto/bcrypt"
	"texApi/config"
	db "texApi/database"
	"texApi/internal/dto"
	"texApi/internal/queries"
)

func GetUser(username, password, loginMethod string) (dto.User, error) {
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
		return user, fmt.Errorf("error fetching user")
	}

	if user.ID == 0 {
		return user, fmt.Errorf("login failed")
	}

	if config.ENV.ENCRYPT_PASSWORDS > 0 {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			return user, fmt.Errorf("login failed")
		}
	} else {
		if user.Password != password {
			return user, fmt.Errorf("login failed")
		}
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
