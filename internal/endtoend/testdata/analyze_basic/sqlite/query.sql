-- name: GetUser :one
SELECT id, name FROM users WHERE id = ?;
