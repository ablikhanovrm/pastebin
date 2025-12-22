-- name: GetUserById :one
SELECT * FROM users AS u WHERE u.id = $1;


-- name: CreateUser :one
INSERT INTO users (username, password_hash)
VALUES ($1, $2)
RETURNING id, username, created_at;