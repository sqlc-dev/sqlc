-- name: GetAuthor :one
SELECT * FROM authors WHERE id = $1;
