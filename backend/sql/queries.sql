-- name: RegisterUser :exec
INSERT INTO users (username, password, email)
VALUES ($1,$2,$3)

-- name: GetUserID :one
SELECT id
FROM users
WHERE username = $1
LIMIT 1;

-- name: ListUsers :many
SELECT *
FROM users
ORDER BY username;

-- name: DeleteUser :exec
DELETE
FROM users
WHERE username = $1;

-- name: LoginUser :exec
UPDATE users
SET logged_in = TRUE
WHERE username = $1;

-- name: CreateNote :exec
INSERT INTO notes (title, user, text, created_at, updated_at)
VALUES ($1,$2,$3,$4,$5)

-- name: GetNote :one
SELECT id
FROM notes
WHERE user = $1 AND title = $2
LIMIT 1;

-- name: DeleteNote :exec
DELETE 
FROM notes
WHERE user = $1 AND title = $2;
