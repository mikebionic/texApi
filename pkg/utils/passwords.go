package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (hash string, err error) {
	var hashByte []byte
	hashByte, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		hash = ""
		return
	}
	hash = string(hashByte)
	return
}

func VerifyPasswordFromHash(password string, hash string) error {
	bHash := []byte(hash)
	return bcrypt.CompareHashAndPassword(bHash, []byte(password))
}
