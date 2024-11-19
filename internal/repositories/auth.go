package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	db "texApi/database"
	"texApi/internal/dto"
	"texApi/internal/queries"
	"time"
)

func GetUser(username, loginMethod string) (dto.User, error) {
	stmt := queries.GetUser
	switch loginMethod {
	case "phone":
		stmt = stmt + " WHERE u.phone = $1"
	case "username":
		stmt = stmt + " WHERE u.username = $1"
	case "email":
		stmt = stmt + " WHERE u.email = $1"
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
	}

	if user.ID == 0 {
		return user, fmt.Errorf("login failed")
	}

	return user, nil
}

func GetUserById(userID int) dto.User {
	var user dto.User
	err := pgxscan.Get(
		context.Background(),
		db.DB,
		&user,
		queries.GetUser+` WHERE u.id = $1`,
		userID,
	)
	if err != nil {
		fmt.Println(err)
	}
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
		user.Phone,
		user.Role,
		user.RoleID,
		user.CompanyID,
		user.Verified,
		user.Active,
		user.OauthProvider,
		user.OauthUserID,
		user.OauthLocation,
		user.OauthAccessToken,
		user.OauthAccessTokenSecret,
		user.OauthRefreshToken,
		user.OauthIDToken,
		user.RefreshToken,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func UpdateUser(user dto.CreateUser, userID int) (int, error) {
	var id int

	err := db.DB.QueryRow(
		context.Background(),
		queries.UpdateUser,
		userID,
		user.Username,
		user.Password,
		user.Email,
		user.Phone,
		user.Role,
		user.RoleID,
		user.CompanyID,
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

func RemoveUserToken(id int) (uID int, err error) {
	err = db.DB.QueryRow(
		context.Background(),
		`UPDATE tbl_user SET refresh_token = '' WHERE id = $1 RETURNING id;`,
		id,
	).Scan(&uID)
	return
}

func ProfileUpdate(user dto.ProfileUpdate, userID int) (int, error) {
	var id int
	err := db.DB.QueryRow(
		context.Background(),
		queries.ProfileUpdate,
		userID,
		user.Username,
		user.Password,
		user.Email,
		user.Phone,
	).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

func SaveUserWithOTP(userID, roleID, verified int, registerType, credentials, otpkey string) (id int, err error) {
	// Update OTP if user exists
	if userID > 0 {
		err = db.DB.QueryRow(
			context.Background(),
			queries.UpdateUserWithOTP,
			registerType,
			credentials,
			roleID,
			verified,
			otpkey,
			userID,
		).Scan(&id)
	} else {
		err = db.DB.QueryRow(
			context.Background(),
			queries.SaveUserWithOTP,
			registerType,
			credentials,
			roleID,
			verified,
			otpkey,
		).Scan(&id)
	}

	if err != nil {
		return 0, err
	}

	return id, nil
}

var ErrOTPExpired = errors.New("otp has expired")
var ErrInvalidOTP = errors.New("invalid otp")
var ErrUserNotFound = errors.New("user not found")

func ValidateOTPAndTime(registerType, credentials, promptOTP string) error {
	var (
		userID     int
		otpKey     string
		verifyTime time.Time
	)

	err := db.DB.QueryRow(
		context.Background(),
		queries.GetOTPInfo,
		registerType,
		credentials,
	).Scan(&userID, &otpKey, &verifyTime)

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrUserNotFound
	} else if err != nil {
		return err
	}

	if otpKey != promptOTP {
		return ErrInvalidOTP
	}

	// Check if the OTP has expired (15 minutes time window)
	expirationTime := verifyTime.Add(15 * time.Minute)
	if time.Now().After(expirationTime) {
		return ErrOTPExpired
	}

	// Set verified = 1
	err = db.DB.QueryRow(
		context.Background(),
		queries.VerifyUserByID,
		userID,
	).Scan(&userID)
	if err != nil {
		return err
	}

	return nil
}
