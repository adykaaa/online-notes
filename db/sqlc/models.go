// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package db

import (
	"database/sql"

	"github.com/google/uuid"
)

type Note struct {
	ID        uuid.UUID
	Title     string
	Username  sql.NullString
	Text      sql.NullString
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
}

type User struct {
	Username string
	Password string
	Email    string
}