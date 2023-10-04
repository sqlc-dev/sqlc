-- name: GetAuthor :one
SELECT * FROM authors AS a
FULL OUTER JOIN books AS b
 ON a.id = b.id
WHERE a.id = ? LIMIT 1;