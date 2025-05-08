package utils

import "golang.org/x/crypto/bcrypt"

func VerifyPassword(passwordByUser string, storedPassword string) bool{
	err := bcrypt.CompareHashAndPassword([]byte(passwordByUser),[]byte(storedPassword))
	if err!=nil{
		return false
	}
	return true
}