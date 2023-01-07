-- name: RegisterUser :exec
INSERT INTO users (email, username, password)
VALUES ($1,$2,$3,$4)

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

-- name: DeleteNote :exec
DELETE 
FROM notes
WHERE id = $1;
