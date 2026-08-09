-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = $1;

-- name: ListAuthors :many
SELECT id, name FROM authors
ORDER BY name;

-- name: CreateAuthor :one
INSERT INTO authors (id, name, bio)
VALUES ($1, $2, $3)
RETURNING *;

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = $1;
