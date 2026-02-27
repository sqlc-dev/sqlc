-- name: GetUserByID :one
SELECT * FROM users WHERE id = ?;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = ?;

-- name: GetActiveUsers :many
SELECT * FROM users WHERE is_active = ?;
