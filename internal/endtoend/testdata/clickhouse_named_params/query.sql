-- name: GetUserByEmail :one
SELECT id, name, email FROM users WHERE email = sqlc.arg('email');

-- name: InsertUser :exec
INSERT INTO users (id, name, email) VALUES (sqlc.arg('id'), sqlc.arg('name'), sqlc.arg('email'));

-- name: FilterUsersByIDs :many
SELECT id, name, email FROM users WHERE id IN (sqlc.slice('ids'));
