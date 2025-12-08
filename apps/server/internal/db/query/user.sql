-- name: GetUserById :one
SELECT * FROM users AS u WHERE u.id = $1;