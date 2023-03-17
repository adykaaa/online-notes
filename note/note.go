package note

import (
	"context"
	"database/sql"
	"errors"
	"time"

	db "github.com/adykaaa/online-notes/db/sqlc"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

var (
	ErrAlreadyExists     = errors.New("note already exists")
	ErrDBInternal        = errors.New("internal DB error during operation")
	ErrNotFound          = errors.New("requested note is not found")
	ErrUserAlreadyExists = errors.New("note already exists")
	ErrUserNotFound      = errors.New("requested note is not found")
)

type service struct {
	q db.Querier
}

func NewService(q db.Querier) *service {
	return &service{q}
}

func (s *service) RegisterUser(ctx context.Context, args *db.RegisterUserParams) (string, error) {
	uname, err := s.q.RegisterUser(ctx, args)

	switch {
	case err != nil:
		if err.(*pq.Error).Code.Name() == "unique_violation" {
			return "", ErrUserAlreadyExists
		}
		return "", ErrDBInternal
	default:
		return uname, nil
	}
}

func (s *service) GetUser(ctx context.Context, username string) (db.User, error) {
	user, err := s.q.GetUser(ctx, username)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return db.User{}, ErrUserNotFound
	case err != nil:
		return db.User{}, ErrDBInternal
	default:
		return user, nil
	}
}

func (s *service) CreateNote(ctx context.Context, title string, username string, text string) (uuid.UUID, error) {
	retID, err := s.q.CreateNote(ctx, &db.CreateNoteParams{
		ID:        uuid.New(),
		Title:     title,
		Username:  username,
		Text:      sql.NullString{String: text, Valid: true},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	switch {
	case err != nil:
		if err.(*pq.Error).Code.Name() == "unique_violation" {
			return uuid.Nil, ErrAlreadyExists
		}
		return uuid.Nil, ErrDBInternal
	default:
		return retID, nil
	}
}

func (s *service) GetAllNotesFromUser(ctx context.Context, username string) ([]db.Note, error) {
	notes, err := s.q.GetAllNotesFromUser(ctx, username)

	if err != nil {
		return nil, ErrDBInternal
	}
	return notes, nil
}

func (s *service) DeleteNote(ctx context.Context, reqID uuid.UUID) (uuid.UUID, error) {
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

func (s *service) UpdateNote(ctx context.Context, reqID uuid.UUID, title string, text string, isTextValid bool) (uuid.UUID, error) {
	id, err := s.q.UpdateNote(ctx, &db.UpdateNoteParams{
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
