package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Note struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	User      string    `json:"user"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (n *Note) Validate() error {
	if n.Title == "" {
		return errors.New("note is missing a title")
	}
	if n.User == "" {
		return errors.New("note is missing the author")
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

	if err := n.Validate(); err != nil {
		return nil, fmt.Errorf("error validating note. %v", err)
	}

	return n, nil
}
