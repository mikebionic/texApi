package repo

import (
	"context"
	"errors"
	"fmt"
	"texApi/config"
	db "texApi/database"
	"texApi/internal/dto"
	"texApi/internal/queries"
	"texApi/pkg/utils"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
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
	default:
		return dto.User{}, fmt.Errorf("invalid login method")
	}

	var users []dto.User
	if err := pgxscan.Select(context.Background(), db.DB, &users, stmt, username); err != nil {
		return dto.User{}, err
	}

	if len(users) == 0 {
		return dto.User{}, fmt.Errorf("user not found")
	}

	return users[0], nil
}

func GetUserById(userID int) (dto.User, error) {
	var users []dto.User
	if err := pgxscan.Select(context.Background(), db.DB, &users, queries.GetUser+" WHERE id = $1", userID); err != nil {
		return dto.User{}, err
	}

	if len(users) == 0 {
		return dto.User{}, fmt.Errorf("user not found")
	}

	return users[0], nil
}

func ManageToken(id int, token, action string) error {
	switch action {
	case "create":
		var users []struct{ ID int }
		err := pgxscan.Select(
			context.Background(),
			db.DB,
			&users,
			`UPDATE tbl_user SET refresh_token = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2 RETURNING id`,
			token,
			id,
		)
		if err != nil {
			return err
		}
		if len(users) == 0 {
			return fmt.Errorf("user not found")
		}
		return nil

	case "validate":
		var tokens []struct{ RefreshToken string }
		err := pgxscan.Select(
			context.Background(),
			db.DB,
			&tokens,
			`SELECT refresh_token FROM tbl_user WHERE id = $1`,
			id,
		)
		if err != nil {
			return err
		}
		if len(tokens) == 0 {
			return fmt.Errorf("user not found")
		}
		if tokens[0].RefreshToken != token {
			return fmt.Errorf("token is invalid")
		}
		return nil

	default:
		return fmt.Errorf("invalid action, must be create or validate")
	}
}

func CreateUser(user dto.CreateUser) (int, error) {
	var result []struct{ ID int }

	if user.OTP == nil {
		user.OTP = &utils.EmptyString
	}

	err := pgxscan.Select(
		context.Background(),
		db.DB,
		&result,
		queries.CreateUser,
		user.Username,
		user.Password,
		user.Email,
		user.Phone,
		user.Role,
		user.RoleID,
		user.CompanyID,
		user.DriverID,
		user.Verified,
		user.Meta,
		user.Meta2,
		user.Meta3,
		user.OTP,
		user.RefreshToken,
		user.Active,
	)

	if err != nil {
		return 0, err
	}
	if len(result) == 0 {
		return 0, fmt.Errorf("user creation failed")
	}
	return result[0].ID, nil
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
		user.DriverID,
		user.Verified,
		user.Active,
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

func UserUpdate(user dto.UserUpdateAuth, userID int) (int, error) {
	var result []struct{ ID int }
	err := pgxscan.Select(
		context.Background(),
		db.DB,
		&result,
		queries.UserUpdate,
		userID,
		user.Username,
		user.Password,
		user.Email,
		user.Phone,
	)

	if err != nil {
		return 0, err
	}
	if len(result) == 0 {
		return 0, fmt.Errorf("user not found")
	}
	return result[0].ID, nil
}

func SaveUserWithOTP(userID, roleID, verified int, registerType, credentials, otp, role string) (int, error) {
	var result []struct{ ID int }
	var err error

	if userID > 0 {
		err = pgxscan.Select(
			context.Background(),
			db.DB,
			&result,
			queries.UpdateUserWithOTP,
			registerType,
			credentials,
			role,
			roleID,
			verified,
			otp,
			userID,
		)
	} else {
		err = pgxscan.Select(
			context.Background(),
			db.DB,
			&result,
			queries.SaveUserWithOTP,
			registerType,
			credentials,
			role,
			roleID,
			verified,
			otp,
		)
	}

	if err != nil {
		return 0, err
	}
	if len(result) == 0 {
		return 0, fmt.Errorf("operation failed")
	}
	return result[0].ID, nil
}

type OTPInfo struct {
	ID         int
	OTPKey     string
	VerifyTime time.Time
}

var ErrOTPExpired = errors.New("otp has expired")
var ErrInvalidOTP = errors.New("invalid otp")
var ErrUserNotFound = errors.New("user not found")

func ValidateOTPAndTime(registerType, credentials, promptOTP string) error {
	var otpInfos []OTPInfo

	err := pgxscan.Select(
		context.Background(),
		db.DB,
		&otpInfos,
		queries.GetOTPInfo,
		registerType,
		credentials,
	)

	if err != nil {
		return err
	}
	if len(otpInfos) == 0 {
		return ErrUserNotFound
	}

	otpInfo := otpInfos[0]
	if otpInfo.OTPKey != promptOTP {
		return ErrInvalidOTP
	}

	expirationTime := otpInfo.VerifyTime.Add(15 * time.Minute)
	if time.Now().Add(config.ENV.TZAddHours).After(expirationTime) {
		return ErrOTPExpired
	}

	// Set verified = 1
	var result []struct{ ID int }
	err = pgxscan.Select(
		context.Background(),
		db.DB,
		&result,
		queries.VerifyUserByID,
		otpInfo.ID,
	)
	if err != nil {
		return err
	}
	if len(result) == 0 {
		return fmt.Errorf("failed to verify user")
	}

	return nil
}

func UpdateUserLastActive(companyID int) error {
	const updateLastActiveQuery = `
		UPDATE tbl_company 
		SET last_active = CURRENT_TIMESTAMP 
		WHERE id = $1
		RETURNING id
	`
	var result []struct{ ID int }
	err := pgxscan.Select(
		context.Background(),
		db.DB,
		&result,
		updateLastActiveQuery,
		companyID,
	)

	if err != nil {
		return fmt.Errorf("failed to update last seen: %w", err)
	}

	if len(result) == 0 {
		return fmt.Errorf("no company found with ID %d", companyID)
	}

	return nil
}
