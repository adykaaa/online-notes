package db

import (
	sqlc "github.com/adykaaa/online-notes/db/sqlc"
)

type Repository struct {
	sqlc.Querier
}

func NewRepository(q sqlc.Querier) *Repository {
	return &Repository{}
}
