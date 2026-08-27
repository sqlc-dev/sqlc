-- name: GetAuthor :one
SELECT *
FROM authors AS a
WHERE p.id = ?
LIMIT 1;
