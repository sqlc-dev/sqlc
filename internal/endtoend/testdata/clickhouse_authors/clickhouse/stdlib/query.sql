-- name: GetAuthor :one
SELECT id, name, bio FROM authors WHERE id = ?;

-- name: ListAuthors :many
SELECT id, name, bio FROM authors ORDER BY name;

-- name: CreateAuthor :exec
INSERT INTO authors (id, name, bio) VALUES (?, ?, ?);
