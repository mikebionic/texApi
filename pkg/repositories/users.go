package repositories

import (
	"context"
	db "texApi/database"
	"texApi/pkg/queries"
	"texApi/pkg/schemas/request"
	"texApi/pkg/schemas/response"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
)

func GetUsers(offset, limit int) ([]response.User, error) {
	var users []response.User

	err := pgxscan.Select(
		context.Background(), db.DB,
		&users, queries.GetUsers, offset, limit,
	)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func GetUser(id int) response.User {
	var user response.User

	db.DB.QueryRow(
		context.Background(), queries.GetUser, id,
	).Scan(
		&user.ID,
		&user.Fullname,
		&user.Phone,
		&user.Address,
		&user.IsVerified,
		&user.SubscriptionID,
	)

	return user
}

func CreateUser(user request.CreateUser) error {
	_, err := db.DB.Exec(
		context.Background(), queries.CreateUser,
		user.Fullname, user.Phone, user.Address, user.Password,
		time.Now(), time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

func CheckSubscription(userID int) int {
	var subscriptionID int

	db.DB.QueryRow(
		context.Background(), queries.CheckSubscription, userID,
	).Scan(&subscriptionID)

	return subscriptionID
}

func BuySubscription(userID, subscriptionID int) error {
	_, err := db.DB.Exec(
		context.Background(), queries.BuySubscription, subscriptionID, userID,
	)

	if err != nil {
		return err
	}

	return nil
}

func CheckUserExist(phone string) string {
	var existPhone string

	db.DB.QueryRow(
		context.Background(), queries.CheckUserExist, phone,
	).Scan(&existPhone)

	return existPhone
}

func CheckUserExistWithStatus(phone string) response.UserExist {
	var user response.UserExist

	db.DB.QueryRow(
		context.Background(), queries.CheckUserExistWithStatus, phone,
	).Scan(&user.Phone, &user.IsVerified)

	return user
}

func UpdateUser(user request.UpdateUser) error {
	if user.Password == "" {
		_, err := db.DB.Exec(
			context.Background(), queries.UpdateUserWithoutPassword,
			user.Fullname, user.Phone, user.Address, user.IsVerified,
			user.SubscriptionID, time.Now(), user.ID,
		)

		if err != nil {
			return err
		}
	} else {
		_, err := db.DB.Exec(
			context.Background(), queries.UpdateUser,
			user.Fullname, user.Phone, user.Address, user.Password,
			user.IsVerified, user.SubscriptionID, time.Now(), user.ID,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateUserPassword(password, phone string) response.Login {
	var user response.Login

	db.DB.QueryRow(
		context.Background(), queries.UpdateUserPassword, password, phone,
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

func DeleteUser(id int) error {
	_, err := db.DB.Exec(context.Background(), queries.DeleteUser, id)

	if err != nil {
		return err
	}

	return nil
}

func VerifyUser(phone string) error {
	_, err := db.DB.Exec(
		context.Background(), queries.VerifyUser, phone,
	)

	if err != nil {
		return err
	}

	return nil
}

func GetUserNotificationToken(userID int) string {
	var token string

	db.DB.QueryRow(
		context.Background(), queries.GetUserNotificationToken, userID,
	).Scan(&token)

	return token
}

func SetUserNotificationToken(token string, userID int) error {
	_, err := db.DB.Exec(
		context.Background(), queries.SetUserNotificationToken, token, userID,
	)

	if err != nil {
		return err
	}

	return nil
}
