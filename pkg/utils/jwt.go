package utils

import (
	"texApi/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(id, roleID, companyID int, role string) (string, string, int64) {
	accessExp := time.Now().Add(config.ENV.ACCESS_TIME).Unix()
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        id,
		"companyID": companyID,
		"role":      role,
		"roleID":    roleID,
		"exp":       accessExp,
	})
	tokenString, _ := accessToken.SignedString([]byte(config.ENV.ACCESS_KEY))

	refreshExp := time.Now().Add(config.ENV.REFRESH_TIME).Unix()
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        id,
		"companyID": companyID,
		"role":      role,
		"roleID":    roleID,
		"exp":       refreshExp,
	})
	refreshString, _ := refreshToken.SignedString([]byte(config.ENV.REFRESH_KEY))

	return tokenString, refreshString, accessExp
}
