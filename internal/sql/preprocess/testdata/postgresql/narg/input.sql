-- name: GetUser :one
SELECT id FROM users WHERE name = sqlc.narg(name);
