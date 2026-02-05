-- name: GetUserById :one
SELECT id, name, email, created_at FROM users AS u WHERE u.id = $1;

-- name: GetUserByEmail :one
SELECT id, name, email, created_at FROM users AS u WHERE u.email = $1;

-- name: CreateUser :one
INSERT INTO users (name, email, password_hash)
VALUES ($1, $2, $3)
RETURNING id, name, created_at;