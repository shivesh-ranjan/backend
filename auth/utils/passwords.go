package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns the Brcypt hash of the password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Can't hash the password: %s", err)
	}
	return string(hashedPassword), nil
}

// CheckPassword crosschecks the hash of the input password with the actual one
func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
