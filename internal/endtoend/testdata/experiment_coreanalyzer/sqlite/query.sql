-- name: GetAuthor :one
SELECT *
FROM authors
WHERE id = ?;

-- name: ListAuthors :many
SELECT id, name
FROM authors
ORDER BY name;

-- name: CreateAuthor :one
INSERT INTO authors (id, name, bio)
VALUES (?, ?, ?)
RETURNING *;

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = ?;
