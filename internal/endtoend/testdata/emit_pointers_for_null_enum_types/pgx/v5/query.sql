-- name: ListUsersByRole :many
SELECT * FROM users WHERE role = $1;

-- name: CreateUser :exec
INSERT INTO users (role, required_role, status) VALUES ($1, $2, $3);
