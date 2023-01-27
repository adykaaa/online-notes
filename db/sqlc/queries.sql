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
WHERE username = $1 LIMIT 1;

-- name: DeleteUser :one
DELETE
FROM users
WHERE username = $1
RETURNING username;

-- name: CreateNote :one
INSERT INTO notes (id,title, username, text, created_at, updated_at)
VALUES ($1,$2,$3,$4,$5,$6)
RETURNING id;

-- name: GetNoteByID :one
SELECT id
FROM notes
WHERE username = $1 AND title = $2;

-- name: GetAllNotesFromUser :many
SELECT *
FROM notes
WHERE username = $1; 

-- name: DeleteNote :one
DELETE 
FROM notes
WHERE username = $1 AND title = $2
RETURNING id;
