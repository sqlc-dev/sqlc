-- name: GetAuthorMv :one
SELECT * FROM authors_mv
WHERE id = $1 LIMIT 1;