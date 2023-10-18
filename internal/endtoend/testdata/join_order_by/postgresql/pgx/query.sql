-- name: GetAuthor :one
SELECT a.name
FROM authors a JOIN authors b ON a.id = b.id
ORDER BY name;
