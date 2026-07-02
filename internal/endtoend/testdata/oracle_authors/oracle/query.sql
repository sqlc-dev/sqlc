-- name: GetAuthor :one
SELECT id, name, bio FROM authors
WHERE id = :1;

-- name: ListAuthors :many
SELECT id, name, bio FROM authors;

-- name: CreateAuthor :exec
INSERT INTO authors (id, name, bio) VALUES (:1, :2, :3);

-- name: DeleteAuthor :exec
DELETE FROM authors WHERE id = :1;
