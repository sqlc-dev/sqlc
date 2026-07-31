-- name: GetUser :one
SELECT id, name FROM users WHERE id = ?;

-- name: FilterUsers :many
SELECT id, name FROM users WHERE name = sqlc.arg(name) AND bio = sqlc.narg(bio);

-- name: CountPosts :one
SELECT count(*) AS total FROM posts WHERE user_id = ?;
