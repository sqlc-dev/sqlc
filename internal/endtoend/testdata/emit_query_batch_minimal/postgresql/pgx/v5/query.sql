-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserCreatedAt :one
SELECT created_at FROM users WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY id;

-- name: UpdateUser :exec
UPDATE users SET name = $1 WHERE id = $2;
