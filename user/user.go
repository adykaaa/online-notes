package user

import (
	"context"

	sqlc "github.com/adykaaa/online-notes/db/sqlc"
)

type UserServicer interface {
	RegisterUser(ctx context.Context, username string, password string, email string) (string, error)
	GetUser(ctx context.Context, username string) (sqlc.User, error)
}
