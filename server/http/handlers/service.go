package handlers

import (
	"context"

	db "github.com/adykaaa/online-notes/db/sqlc"
	"github.com/google/uuid"
)

type NoteService interface {
	CreateNote(ctx context.Context, title string, username string, text string) (uuid.UUID, error)
	GetAllNotesFromUser(ctx context.Context, username string) ([]db.Note, error)
	DeleteNote(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
	UpdateNote(ctx context.Context, reqID uuid.UUID, title string, text string, isTextEmpty bool) (uuid.UUID, error)
	RegisterUser(ctx context.Context, args *db.RegisterUserParams) (string, error)
	GetUser(ctx context.Context, username string) (db.User, error)
}
