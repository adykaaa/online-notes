package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Note struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	Title     string `json:"title"`
	User      string `json:"user"`
	Text      string `json:"text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (n *Note) Validate() error {
	if n.Title == "" {
		return errors.New("Note is missing a title!")
	}
	if n.User == "" {
		return errors.New("Note is missing the author!")
	}

	return nil
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

	err := n.Validate()
	if err != nil {
		return nil, fmt.Errorf("Error validating node! %v", err)
	}
	return n, nil
}
