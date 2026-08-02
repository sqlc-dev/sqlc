-- name: GetAuthor :one
SELECT * FROM authors WHERE id = $1 LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM authors ORDER BY name;

-- name: CreateAuthor :one
INSERT INTO authors (name, bio, status) VALUES ($1, $2, $3) RETURNING *;

-- name: ListAuthorsByStatus :many
SELECT * FROM authors WHERE status = $1 ORDER BY name;

-- name: GetBook :one
SELECT * FROM books WHERE id = $1 LIMIT 1;
