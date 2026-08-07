-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = ?;

-- name: ListAuthors :many
SELECT id, name FROM authors
ORDER BY name;

-- name: CreateAuthor :execresult
INSERT INTO authors (id, name, bio)
VALUES (?, ?, ?);

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = ?;
