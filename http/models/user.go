package http

import (
	"errors"
	"fmt"
)

type User struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	LoggedIn bool   `json:"loggedIn"`
}

func (u *User) ValidateUser() error {
	if u.Email == "" {
		return errors.New("user email cannot be nil")
	}

	if u.Username == "" {
		return errors.New("username cannot be nil")
	}

	if u.Password == "" {
		return errors.New("password cannot be nil")
	}
	return nil
}

func NewUser(email string, username string, password string, loggedIn bool) (*User, error) {
	u := &User{
		Email:    email,
		Username: username,
		Password: password,
		LoggedIn: loggedIn,
	}
	if err := u.ValidateUser(); err != nil {
		return nil, fmt.Errorf("error validating the User! %v", err)
	}
	return u, nil
}
