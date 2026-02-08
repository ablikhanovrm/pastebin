-- name: GetPasteById :one
SELECT * FROM pastes as p
WHERE p.id = $1;

-- name: GetPastes :many
SELECT * FROM pastes as p
WHERE p.user_id = $1
ORDER BY p.created_at DESC;

-- name: CreatePaste :one
INSERT INTO pastes (uuid, user_id, title, content, syntax, is_private, is_burn_after_read, expire_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;