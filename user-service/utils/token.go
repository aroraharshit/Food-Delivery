package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateToken(mobile string) (string, error) {
	JWT_SECRET_KEY := os.Getenv("JWT_SECRET_KEY")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": mobile,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(JWT_SECRET_KEY)
	if err != nil {
		return "", fmt.Errorf("Error signing the token: %v", err)
	}

	return tokenString, nil
}
