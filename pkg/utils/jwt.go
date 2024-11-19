package utils

import (
	"texApi/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(id, roleID, companyID int, role string) (string, string, int64) {
	unixTime := time.Now().Add(config.ENV.REFRESH_TIME).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        id,
		"roleID":    roleID,
		"companyID": companyID,
		"role":      role,
		"exp":       unixTime,
	})

	tokenString, _ := token.SignedString([]byte(config.ENV.ACCESS_KEY))
	refreshString, _ := token.SignedString([]byte(config.ENV.REFRESH_KEY))

	return tokenString, refreshString, unixTime
}
