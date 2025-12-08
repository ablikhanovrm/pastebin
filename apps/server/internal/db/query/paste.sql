-- name: GetPasteById :one
SELECT * FROM pastes as p WHERE p.id = $1;