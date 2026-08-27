-- name: GetAuthor :one
SELECT *
FROM authors AS a
FULL JOIN books AS b ON a.id = b.id
WHERE a.id = ?
LIMIT 1;
