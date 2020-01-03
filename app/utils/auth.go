package utils

import (
	"errors"
	"strings"
)

const (
	HeaderBearerType  = "Bearer"
)

func GetTokenFromAuthHeader(authHeader string, headerType string) (string, error) {
	splitHeader := strings.Split(authHeader, headerType + " ")
	if len(splitHeader) != 2 {
		return "", errors.New("header not valid")
	}
	return splitHeader[1], nil
}