-- name: ByIDs :many
SELECT id, name FROM items WHERE id IN (sqlc.slice(ids));
