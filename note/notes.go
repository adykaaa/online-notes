package note

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sqlc "github.com/adykaaa/online-notes/db/sqlc"
	"github.com/google/uuid"
)

type NoteServicer interface {
	CreateNote(ctx context.Context, title string, username string, text string) (uuid.UUID, error)
	GetAllNotesFromUser(ctx context.Context, username string) ([]sqlc.Note, error)
	DeleteNote(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
	UpdateNote(ctx context.Context, reqID uuid.UUID, title string, text string, isTextEmpty bool) (sqlc.Note, error)
}

var (
	ErrAlreadyExists = errors.New("note already exists")
	ErrDBInternal    = errors.New("internal DB error during operation")
	ErrNotFound      = errors.New("requested Note is not found")
)

type NoteService struct {
	q sqlc.Querier
}

func (ns *NoteService) CreateNote(ctx context.Context, title string, username string, text string) (uuid.UUID, error) {
	retID, err := ns.q.CreateNote(ctx, &sqlc.CreateNoteParams{
		ID:        uuid.New(),
		Title:     title,
		Username:  username,
		Text:      sql.NullString{String: text, Valid: true},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return uuid.Nil, ErrNotFound
	case err != nil:
		return uuid.Nil, ErrDBInternal
	default:
		return retID, nil
	}
}

func (ns *NoteService) GetAllNotesFromUser(ctx context.Context, username string) ([]sqlc.Note, error) {
	notes, err := ns.q.GetAllNotesFromUser(ctx, username)

	if err != nil {
		return nil, ErrDBInternal
	}
	return notes, nil
}

func (ns *NoteService) DeleteNote(ctx context.Context, reqID uuid.UUID) (uuid.UUID, error) {
	id, err := ns.q.DeleteNote(ctx, reqID)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return uuid.Nil, ErrNotFound
	case err != nil:
		return uuid.Nil, ErrDBInternal
	default:
		return id, nil
	}
}

func (ns *NoteService) UpdateNote(ctx context.Context, reqID uuid.UUID, title string, text string, isTextEmpty bool) (uuid.UUID, error) {
	id, err := ns.q.UpdateNote(ctx, &sqlc.UpdateNoteParams{
		ID:        reqID,
		Title:     sql.NullString{String: title, Valid: true},
		Text:      sql.NullString{String: text, Valid: isTextEmpty},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	if err != nil {
		return uuid.Nil, ErrDBInternal
	}
	return id, nil
}
