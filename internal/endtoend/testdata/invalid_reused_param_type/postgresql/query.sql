-- name: CreateAuthor :one
INSERT INTO authors (
	name, bio, age
) VALUES (
	$1, $1, $1
)
RETURNING *;