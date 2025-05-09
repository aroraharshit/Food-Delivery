package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func VerifyPassword(passwordByUser string, storedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(passwordByUser))
	return err == nil
}

func HashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}
