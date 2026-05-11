-- name: CreateAuthor :one
INSERT INTO authors (name, status) VALUES ($1, $2) RETURNING *;

-- name: GetAuthor :one
SELECT * FROM authors WHERE id = $1 LIMIT 1;
