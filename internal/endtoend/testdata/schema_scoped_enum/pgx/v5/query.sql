-- name: ListUsersByRole :many
SELECT * FROM foo.users WHERE role = $1;
