package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (hashed string, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	hashed = string(hash)
	return
}

func VerifyPassword(password, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)) == nil
}
