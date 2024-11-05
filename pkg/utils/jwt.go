package utils

import (
	"texApi/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(id, roleID int) (string, int64) {
	unixTime := time.Now().Add(config.ENV.REFRESH_TIME).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":     id,
		"roleID": roleID,
		"exp":    unixTime,
	})

	tokenString, _ := token.SignedString([]byte(config.ENV.ACCESS_KEY))

	return tokenString, unixTime
}
