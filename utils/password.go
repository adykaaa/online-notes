package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashUserPassword(password string) (string, error) {
	if len(password) < 5 {
		return "", fmt.Errorf("password cannot be less than 5 chars long")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", fmt.Errorf("could not hash user password! %v", err)
	}
	return string(hash), nil
}

func ValidateUserPassword(hashedPassword string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("error while comparing the current and hashed password! %v", err)
	}
	return nil
}
