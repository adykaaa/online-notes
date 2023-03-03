-- name: RegisterUser :one
INSERT INTO users (username, password, email)
VALUES ($1,$2,$3)
RETURNING username;

-- name: ListUsers :many
SELECT *
FROM users
ORDER BY username;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1;

-- name: CreateNote :one
INSERT INTO notes (id, title, username, text, created_at, updated_at)
VALUES ($1,$2,$3,$4,$5,$6)
RETURNING id;

-- name: UpdateNote :one
UPDATE notes
SET
  title = COALESCE(sqlc.narg(title), title),
  text = COALESCE(sqlc.narg(text), text),
  updated_at = COALESCE(sqlc.narg(updated_at), updated_at)
WHERE
  id = $1
RETURNING id;

-- name: GetAllNotesFromUser :many
SELECT *
FROM notes
WHERE username = $1; 

-- name: DeleteNote :one
DELETE 
FROM notes
WHERE id = $1
RETURNING id;
