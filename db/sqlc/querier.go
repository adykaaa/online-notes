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
	CreateNote(ctx context.Context, arg *CreateNoteParams) (Note, error)
	DeleteNote(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
	DeleteUser(ctx context.Context, username string) (string, error)
	GetAllNotesFromUser(ctx context.Context, username sql.NullString) ([]Note, error)
	GetNoteByID(ctx context.Context, arg *GetNoteByIDParams) (uuid.UUID, error)
	GetUser(ctx context.Context, username string) (User, error)
	ListUsers(ctx context.Context) ([]User, error)
	RegisterUser(ctx context.Context, arg *RegisterUserParams) (string, error)
	UpdateNote(ctx context.Context, arg *UpdateNoteParams) (Note, error)
}

var _ Querier = (*Queries)(nil)
