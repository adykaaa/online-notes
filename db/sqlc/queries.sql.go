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

const createNote = `-- name: CreateNote :one
INSERT INTO notes (id,title, username, text, created_at, updated_at)
VALUES ($1,$2,$3,$4,$5,$6)
RETURNING id, title, username, text, created_at, updated_at
`

type CreateNoteParams struct {
	ID        uuid.UUID
	Title     string
	Username  sql.NullString
	Text      sql.NullString
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
}

func (q *Queries) CreateNote(ctx context.Context, arg CreateNoteParams) (Note, error) {
	row := q.db.QueryRowContext(ctx, createNote,
		arg.ID,
		arg.Title,
		arg.Username,
		arg.Text,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Note
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Username,
		&i.Text,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteNote = `-- name: DeleteNote :one
DELETE 
FROM notes
WHERE username = $1 AND title = $2
RETURNING id
`

type DeleteNoteParams struct {
	Username sql.NullString
	Title    string
}

func (q *Queries) DeleteNote(ctx context.Context, arg DeleteNoteParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, deleteNote, arg.Username, arg.Title)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const deleteUser = `-- name: DeleteUser :one
DELETE
FROM users
WHERE username = $1
RETURNING username
`

func (q *Queries) DeleteUser(ctx context.Context, username string) (string, error) {
	row := q.db.QueryRowContext(ctx, deleteUser, username)
	err := row.Scan(&username)
	return username, err
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
	items := []Note{}
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

const getUser = `-- name: GetUser :one
SELECT username, password, email FROM users
WHERE username = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, username)
	var i User
	err := row.Scan(&i.Username, &i.Password, &i.Email)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT username, password, email
FROM users
ORDER BY username
`

func (q *Queries) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(&i.Username, &i.Password, &i.Email); err != nil {
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

const registerUser = `-- name: RegisterUser :one
INSERT INTO users (username, password, email)
VALUES ($1,$2,$3)
RETURNING username
`

type RegisterUserParams struct {
	Username string
	Password string
	Email    string
}

func (q *Queries) RegisterUser(ctx context.Context, arg RegisterUserParams) (string, error) {
	row := q.db.QueryRowContext(ctx, registerUser, arg.Username, arg.Password, arg.Email)
	var username string
	err := row.Scan(&username)
	return username, err
}

const updateNoteText = `-- name: UpdateNoteText :one
UPDATE notes SET text = $1, updated_at = $2 WHERE id = $3 RETURNING id, title, username, text, created_at, updated_at
`

type UpdateNoteTextParams struct {
	Text      sql.NullString
	UpdatedAt sql.NullTime
	ID        uuid.UUID
}

func (q *Queries) UpdateNoteText(ctx context.Context, arg UpdateNoteTextParams) (Note, error) {
	row := q.db.QueryRowContext(ctx, updateNoteText, arg.Text, arg.UpdatedAt, arg.ID)
	var i Note
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Username,
		&i.Text,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateNoteTitle = `-- name: UpdateNoteTitle :one
UPDATE notes SET title = $1, updated_at = $2 WHERE id = $3 RETURNING id, title, username, text, created_at, updated_at
`

type UpdateNoteTitleParams struct {
	Title     string
	UpdatedAt sql.NullTime
	ID        uuid.UUID
}

func (q *Queries) UpdateNoteTitle(ctx context.Context, arg UpdateNoteTitleParams) (Note, error) {
	row := q.db.QueryRowContext(ctx, updateNoteTitle, arg.Title, arg.UpdatedAt, arg.ID)
	var i Note
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Username,
		&i.Text,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
