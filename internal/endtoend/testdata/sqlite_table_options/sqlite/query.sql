-- name: GetAuthor :one
SELECT * FROM authors1
WHERE id = ?1 LIMIT 1;
