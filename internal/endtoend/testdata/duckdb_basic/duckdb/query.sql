-- name: GetAuthor :one
SELECT id, name, bio FROM authors
WHERE id = $1;

-- name: ListAuthors :many
SELECT id, name, bio FROM authors
ORDER BY name;

-- name: CreateAuthor :one
INSERT INTO authors (name, bio)
VALUES ($1, $2)
RETURNING id, name, bio;

-- name: DeleteAuthor :exec
DELETE FROM authors WHERE id = $1;
