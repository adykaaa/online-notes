package models

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
	Title     string    `json:"title" validate:"required,min=4"`
	User      string    `json:"user" validate:"required"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
