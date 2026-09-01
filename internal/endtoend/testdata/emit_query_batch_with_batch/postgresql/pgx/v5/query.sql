-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY id;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: BatchUpdateUser :batchexec
UPDATE users SET name = $1 WHERE id = $2;
