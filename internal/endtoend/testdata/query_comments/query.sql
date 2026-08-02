-- name: GetAuthor :one
SELECT id, name FROM authors WHERE id = $1;

-- name: ListAuthors :many
SELECT id, name FROM authors ORDER BY name;
