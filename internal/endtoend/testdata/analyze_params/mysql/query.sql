-- name: GetUser :one
SELECT id, name FROM users WHERE id = ?;

-- name: FilterUsers :many
SELECT id, name FROM users WHERE name = sqlc.arg(name) AND bio = sqlc.narg(bio);

-- name: UsersByID :many
SELECT id, name FROM users WHERE id IN (sqlc.slice(ids));
