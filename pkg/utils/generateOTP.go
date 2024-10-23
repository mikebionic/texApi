package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateOTP(length ...int) (string, error) {
	otpLength := 4
	if len(length) > 0 && length[0] > 0 {
		otpLength = length[0]
	}

	otp := ""
	for i := 0; i < otpLength; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		otp += fmt.Sprintf("%d", num)
	}
	return otp, nil
}
