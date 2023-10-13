-- name: CreateAuthor :one
INSERT INTO authors (
  name, bio, missing_column
) VALUES (
  $1, $2, true
)
RETURNING *;
