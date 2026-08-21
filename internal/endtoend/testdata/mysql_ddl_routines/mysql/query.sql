-- name: GetCollection :one
SELECT id, updated_at FROM collections WHERE id = ?;
