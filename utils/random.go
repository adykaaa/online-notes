package utils

import (
	"database/sql"
	"math/rand"
	"strings"
	"time"

	db "github.com/adykaaa/online-notes/db/sqlc"
	"github.com/google/uuid"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func NewRandomString(n int) string {
	var sb strings.Builder
	k := len(chars)

	for i := 0; i < n; i++ {
		c := chars[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func NewRandomDBNote(id uuid.UUID) *db.Note {
	note := db.Note{
		ID:        id,
		Title:     NewRandomString(15),
		Username:  NewRandomString(10),
		Text:      sql.NullString{String: NewRandomString(60), Valid: true},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return &note
}
