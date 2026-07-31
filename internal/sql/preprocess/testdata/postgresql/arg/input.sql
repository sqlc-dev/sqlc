-- name: GetUser :one
SELECT id FROM users WHERE name = sqlc.arg(name);
