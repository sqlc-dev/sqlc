-- name: ListUsers :many
SELECT name FROM users WHERE (name = sqlc.narg(name) OR sqlc.narg(name) IS NULL) AND age > sqlc.arg(min_age);
