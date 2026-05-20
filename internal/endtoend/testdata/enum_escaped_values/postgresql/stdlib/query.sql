-- name: GetUser :one
SELECT id, role, payload FROM users WHERE id = $1;
