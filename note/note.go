package note

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sqlc "github.com/adykaaa/online-notes/db/sqlc"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Servicer interface {
	CreateNote(ctx context.Context, title string, username string, text string) (uuid.UUID, error)
	GetAllNotesFromUser(ctx context.Context, username string) ([]sqlc.Note, error)
	DeleteNote(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
	UpdateNote(ctx context.Context, reqID uuid.UUID, title string, text string, isTextEmpty bool) (sqlc.Note, error)
}

var (
	ErrAlreadyExists = errors.New("note already exists")
	ErrDBInternal    = errors.New("internal DB error during operation")
	ErrNotFound      = errors.New("requested note is not found")
)

type Service struct {
	q sqlc.Querier
}

func (s *Service) CreateNote(ctx context.Context, title string, username string, text string) (uuid.UUID, error) {
	retID, err := s.q.CreateNote(ctx, &sqlc.CreateNoteParams{
		ID:        uuid.New(),
		Title:     title,
		Username:  username,
		Text:      sql.NullString{String: text, Valid: true},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	switch {
	case err.(*pq.Error).Code.Name() == "unique_violation":
		return uuid.Nil, ErrAlreadyExists
	case err != nil:
		return uuid.Nil, ErrDBInternal
	default:
		return retID, nil
	}
}

func (s *Service) GetAllNotesFromUser(ctx context.Context, username string) ([]sqlc.Note, error) {
	notes, err := s.q.GetAllNotesFromUser(ctx, username)

	if err != nil {
		return nil, ErrDBInternal
	}
	return notes, nil
}

func (s *Service) DeleteNote(ctx context.Context, reqID uuid.UUID) (uuid.UUID, error) {
	id, err := s.q.DeleteNote(ctx, reqID)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return uuid.Nil, ErrNotFound
	case err != nil:
		return uuid.Nil, ErrDBInternal
	default:
		return id, nil
	}
}

func (s *Service) UpdateNote(ctx context.Context, reqID uuid.UUID, title string, text string, isTextValid bool) (uuid.UUID, error) {
	id, err := s.q.UpdateNote(ctx, &sqlc.UpdateNoteParams{
		ID:        reqID,
		Title:     sql.NullString{String: title, Valid: true},
		Text:      sql.NullString{String: text, Valid: isTextValid},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return uuid.Nil, ErrNotFound
	case err != nil:
		return uuid.Nil, ErrDBInternal
	default:
		return id, nil
	}
}
