-- name: GetUser :one
SELECT id, name FROM users WHERE id = $1;

-- name: FilterUsers :many
SELECT id, name FROM users WHERE name = sqlc.arg(name) AND bio = sqlc.narg(bio);

-- name: UsersByID :many
SELECT id, name FROM users WHERE id = ANY(sqlc.slice(ids));

-- name: PostsBetween :many
SELECT id, title FROM posts WHERE created BETWEEN $1 AND $2;
