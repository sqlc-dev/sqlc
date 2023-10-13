-- name: GetAuthor :one
SELECT * FROM authors
WHERE name = $1 AND country_code = $2 LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM authors
ORDER BY name;

-- name: CreateAuthor :one
INSERT INTO authors (
  name, bio, country_code
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = $1;
