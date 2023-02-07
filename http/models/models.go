package http

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Username string `json:"username" validate:"required,min=5,max=30,alphanum"`
	Password string `json:"password" validate:"required,min=5"`
	Email    string `json:"email" validate:"email,required"`
}

type Note struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title" validate:"required,min=4,alphanum"`
	User      string    `json:"user" validate:"required"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewUser(email string, username string, password string) (*User, error) {
	u := &User{
		Email:    email,
		Username: username,
		Password: password,
	}
	return u, nil
}

func NewNote(title string, text string, user string) (*Note, error) {
	n := &Note{
		ID:        uuid.New(),
		Title:     title,
		User:      user,
		Text:      text,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return n, nil
}
