package password

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var ErrTooShort = errors.New("the given password is too short")

func Hash(password string) (string, error) {
	if len(password) < 5 {
		return "", ErrTooShort
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", fmt.Errorf("could not hash user password! %v", err)
	}
	return string(hash), nil
}

func Validate(hashedPassword string, password string) error {
	if len(password) < 5 {
		return fmt.Errorf("password cannot be less than 5 chars long")
	}
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("error while comparing the current and hashed password! %v", err)
	}
	return nil
}
