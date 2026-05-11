-- name: ListUsersByRole :many
SELECT * FROM users WHERE role = $1;
