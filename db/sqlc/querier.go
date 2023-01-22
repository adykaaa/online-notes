// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Querier interface {
	CreateNote(ctx context.Context, arg CreateNoteParams) error
	DeleteNote(ctx context.Context, arg DeleteNoteParams) error
	DeleteUser(ctx context.Context, username string) error
	GetAllNotesFromUser(ctx context.Context, username sql.NullString) ([]Note, error)
	GetNoteByID(ctx context.Context, arg GetNoteByIDParams) (uuid.UUID, error)
	ListUsers(ctx context.Context) ([]User, error)
	RegisterUser(ctx context.Context, arg RegisterUserParams) error
}

var _ Querier = (*Queries)(nil)
