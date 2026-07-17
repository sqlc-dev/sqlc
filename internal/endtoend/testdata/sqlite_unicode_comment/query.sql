-- name: GetItem :one
SELECT id, name FROM items WHERE id = ?;

-- section — divider

-- name: UpdateItem :exec
UPDATE items SET name = ? WHERE id = ?;
