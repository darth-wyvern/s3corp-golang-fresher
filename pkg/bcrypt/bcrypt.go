package bcrypt

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hash password using the provided password and the provided algorithm
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash verify password using the provided password and the provided algorithm
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
