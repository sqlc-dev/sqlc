-- name: GetAuthor :one
SELECT * FROM authors WHERE id = ?;

-- name: ListAuthors :many
SELECT * FROM authors;

-- name: CreateAuthor :one
INSERT INTO authors (name, bio) VALUES (?, ?) RETURNING *;
