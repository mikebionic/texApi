package utils

import (
	"encoding/base64"
	"errors"
	"strings"
)

func ExtractBasicAuth(authHeader string) (string, string, error) {
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "basic" {
		return "", "", errors.New("invalid authorization header format")
	}

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", errors.New("failed to decode base64 string")
	}

	credentials := strings.SplitN(string(decoded), ":", 2)
	if len(credentials) != 2 {
		return "", "", errors.New("invalid basic auth credentials")
	}

	return credentials[0], credentials[1], nil
}
