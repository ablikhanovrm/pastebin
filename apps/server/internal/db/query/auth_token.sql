-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id, token_hash, user_agent, ip_address, expires_at) VALUES ($1,$2,$3,$4,$5) RETURNING id, created_at;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens 
SET revoked = true, revoked_at = NOW()
WHERE token_hash = $1 AND revoked_at IS NULL;

-- name: GetRefreshTokenByHash :one
SELECT r.user_id, r.token_hash, r.user_agent, r.ip_address, r.expires_at 
FROM refresh_tokens AS r
WHERE r.token_hash = $1 AND r.revoked = false
LIMIT 1;