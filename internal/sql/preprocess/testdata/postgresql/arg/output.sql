-- name: GetUser :one
SELECT id FROM users WHERE name = $1;
