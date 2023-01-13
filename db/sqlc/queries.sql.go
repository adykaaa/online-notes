// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: queries.sql

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createNote = `-- name: CreateNote :exec
INSERT INTO notes (title, username, text, created_at, updated_at)
VALUES ($1,$2,$3,$4,$5)
`

type CreateNoteParams struct {
	Title     string
	Username  sql.NullString
	Text      sql.NullString
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
}

func (q *Queries) CreateNote(ctx context.Context, arg CreateNoteParams) error {
	_, err := q.db.ExecContext(ctx, createNote,
		arg.Title,
		arg.Username,
		arg.Text,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}

const deleteNote = `-- name: DeleteNote :exec
DELETE 
FROM notes
WHERE username = $1 AND title = $2
`

type DeleteNoteParams struct {
	Username sql.NullString
	Title    string
}

func (q *Queries) DeleteNote(ctx context.Context, arg DeleteNoteParams) error {
	_, err := q.db.ExecContext(ctx, deleteNote, arg.Username, arg.Title)
	return err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE
FROM users
WHERE username = $1
`

func (q *Queries) DeleteUser(ctx context.Context, username string) error {
	_, err := q.db.ExecContext(ctx, deleteUser, username)
	return err
}

const getAllNotesFromUser = `-- name: GetAllNotesFromUser :many
SELECT id, title, username, text, created_at, updated_at
FROM notes
WHERE username = $1
`

func (q *Queries) GetAllNotesFromUser(ctx context.Context, username sql.NullString) ([]Note, error) {
	rows, err := q.db.QueryContext(ctx, getAllNotesFromUser, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Note
	for rows.Next() {
		var i Note
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Username,
			&i.Text,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNoteByID = `-- name: GetNoteByID :one
SELECT id
FROM notes
WHERE username = $1 AND title = $2
`

type GetNoteByIDParams struct {
	Username sql.NullString
	Title    string
}

func (q *Queries) GetNoteByID(ctx context.Context, arg GetNoteByIDParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, getNoteByID, arg.Username, arg.Title)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const listUsers = `-- name: ListUsers :many
SELECT username, password, email, logged_in
FROM users
ORDER BY username
`

func (q *Queries) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.Username,
			&i.Password,
			&i.Email,
			&i.LoggedIn,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const loginUser = `-- name: LoginUser :exec
UPDATE users
SET logged_in = TRUE
WHERE username = $1
`

func (q *Queries) LoginUser(ctx context.Context, username string) error {
	_, err := q.db.ExecContext(ctx, loginUser, username)
	return err
}

const registerUser = `-- name: RegisterUser :exec
INSERT INTO users (username, password, email)
VALUES ($1,$2,$3)
`

type RegisterUserParams struct {
	Username string
	Password string
	Email    string
}

func (q *Queries) RegisterUser(ctx context.Context, arg RegisterUserParams) error {
	_, err := q.db.ExecContext(ctx, registerUser, arg.Username, arg.Password, arg.Email)
	return err
}