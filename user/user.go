package user

import (
	"context"
	"database/sql"
	"errors"

	sqlc "github.com/adykaaa/online-notes/db/sqlc"
	"github.com/lib/pq"
)

var (
	ErrAlreadyExists = errors.New("user already exists")
	ErrDBInternal    = errors.New("internal DB error during operation")
	ErrNotFound      = errors.New("requested user is not found")
)

type Servicer interface {
	RegisterUser(ctx context.Context, username string, password string, email string) (string, error)
	GetUser(ctx context.Context, username string) (sqlc.User, error)
}

type Service struct {
	q sqlc.Querier
}

func (s *Service) RegisterUser(ctx context.Context, username string, hashedpw string, email string) (string, error) {
	uname, err := s.q.RegisterUser(ctx, &sqlc.RegisterUserParams{
		Username: username,
		Password: hashedpw,
		Email:    email,
	})

	switch {
	case err.(*pq.Error).Code.Name() == "unique_violation":
		return "", ErrAlreadyExists
	case err != nil:
		return "", ErrDBInternal
	default:
		return uname, nil
	}
}

func (s *Service) GetUser(ctx context.Context, username string) (sqlc.User, error) {
	user, err := s.q.GetUser(ctx, username)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return sqlc.User{}, ErrNotFound
	case err != nil:
		return sqlc.User{}, ErrDBInternal
	default:
		return user, nil
	}
}
