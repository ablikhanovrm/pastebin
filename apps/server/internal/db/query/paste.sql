-- name: GetPasteById :one
SELECT * FROM pastes as p
WHERE p.uuid = $1 AND p.user_id = $2
LIMIT 1;

-- name: GetPastesFirstPage :many
SELECT * FROM pastes as p
WHERE p.visibility='public'
   OR p.user_id=$1 AND (expires_at IS NULL OR expires_at > now())
ORDER BY created_at DESC
LIMIT $2;

-- name: GetPastesAfterCursor :many
SELECT * FROM pastes as p
WHERE p.visibility='public'
   OR p.user_id=$1 AND (expires_at IS NULL OR expires_at > now()) AND created_at < $2
ORDER BY p.created_at DESC
LIMIT $3;


-- name: GetUserPastesFirstPage :many
SELECT * FROM pastes as p
WHERE p.user_id=$1 AND (expires_at IS NULL OR expires_at > now())
ORDER BY created_at DESC
LIMIT $2;

-- name: GetUserPastesAfterCursor :many
SELECT * FROM pastes as p
WHERE p.user_id=$1 AND (expires_at IS NULL OR expires_at > now()) AND created_at < $2
ORDER BY p.created_at DESC
LIMIT $3;


-- name: CreatePaste :one
INSERT INTO pastes (uuid, user_id, title, s3_key, syntax, max_views, visibility, expire_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdatePaste :one
UPDATE pastes
SET
    title = $3,
    syntax = $4,
    visibility = $5,
    max_views = $6,
    expire_at = $7,
    updated_at = now()
WHERE uuid = $1 AND user_id = $2
RETURNING *;

-- name: DeletePaste :exec
DELETE FROM pastes WHERE uuid = $1 AND user_id = $2;