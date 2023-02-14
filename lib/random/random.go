package random

import (
	"database/sql"
	"math/rand"
	"strings"
	"time"

	db "github.com/adykaaa/online-notes/db/sqlc"
	"github.com/google/uuid"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func NewString(n int) string {
	var sb strings.Builder
	k := len(chars)

	for i := 0; i < n; i++ {
		c := chars[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func NewDBNote(id uuid.UUID) *db.Note {
	note := db.Note{
		ID:        id,
		Title:     NewString(15),
		Username:  NewString(10),
		Text:      sql.NullString{String: NewString(60), Valid: true},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return &note
}
