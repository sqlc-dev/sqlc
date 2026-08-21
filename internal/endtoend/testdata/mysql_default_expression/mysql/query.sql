-- name: GetExpiration :one
SELECT id, created_at, expiry FROM expirations WHERE id = ?;
