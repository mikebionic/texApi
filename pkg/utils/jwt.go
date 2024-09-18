package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(
	id int, expiration time.Duration, secret_key, role string,
) string {
	unixTime := time.Now().Add(expiration).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   id,
		"role": role,
		"exp":  unixTime,
	})

	tokenString, _ := token.SignedString([]byte(secret_key))

	return tokenString
}
