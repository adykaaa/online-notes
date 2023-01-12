package models

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
	Password string    `json:"password"`
}

func (u *User) ValidateUser() error {
	if u.Email == "" {
		return errors.New("User email cannot be nil!")
	}

	if u.Username == "" {
		return errors.New("Username cannot be nil!")
	}

	if u.Password == "" {
		return errors.New("Password cannot be nil!")
	}
	return nil
}

func NewUser(email string, username string, password string) (*User, error) {
	u := &User{
		ID:       uuid.New(),
		Email:    email,
		Username: username,
		Password: password,
	}
	if err := u.ValidateUser(); err != nil {
		return nil, fmt.Errorf("Error validating the User! %v", err)
	}
	return u, nil
}

func HashUserPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", fmt.Errorf("Could not hash user password! %v", err)
	}
	return string(hash), nil
}

func (u *User) ValidateUserPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return fmt.Errorf("Error while comparing the current and hashed password! %v", err)
	}
	return nil
}
