-- name: ListAuthors :many
SELECT * FROM authors;

-- name: GetAuthor :one
SELECT * FROM authors WHERE id = $1;

-- name: CreateAuthor :one
INSERT INTO authors (name, bio) VALUES ($1, $2) RETURNING *;

-- name: UpdateAuthor :one
UPDATE authors SET name = $1, bio = $2 WHERE id = $3 RETURNING *;

-- name: DeleteAuthor :one
DELETE FROM authors WHERE id = $1 RETURNING *;
