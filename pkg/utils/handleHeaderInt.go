package utils

import (
	"strconv"
)

func HandleHeaderInt(headerValue string) (intVal int, err error) {
	if headerValue != "" {
		if intVal, err := strconv.Atoi(headerValue); err == nil {
			return intVal, nil
		}
	}
	return
}
