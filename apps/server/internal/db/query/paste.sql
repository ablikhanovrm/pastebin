-- name: GetPasteById :one
SELECT * FROM pastes as p
WHERE p.uuid = $1
LIMIT 1;

-- name: GetPastes :many
SELECT * FROM pastes as p
WHERE p.visibility='public'
   OR p.user_id=$1
ORDER BY p.created_at DESC;

-- name: GetUserPastes :many
SELECT * FROM pastes as p
WHERE p.user_id = $1
ORDER BY p.created_at DESC;

-- name: CreatePaste :one
INSERT INTO pastes (uuid, user_id, title, s3_key, syntax, max_views, visibility, expire_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdatePaste :exec
UPDATE pastes
SET
    title = $2,
    syntax = $3,
    visibility = $4,
    max_views = $5,
    expire_at = $6,
    updated_at = now()
WHERE uuid = $1;

-- name: DeletePaste :exec
DELETE FROM pastes WHERE uuid = $1;