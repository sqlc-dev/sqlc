-- name: GetAuthor :one
SELECT  *
FROM    authors a
WHERE   p.id = $1
LIMIT   1;
