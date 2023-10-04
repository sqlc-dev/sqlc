-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = $1 LIMIT 1;
