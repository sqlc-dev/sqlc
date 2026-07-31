-- name: GetEvent :one
SELECT id, name FROM events WHERE id = ?;

-- name: FilterEvents :many
SELECT id, name FROM events WHERE name = ? AND amount > ?;
