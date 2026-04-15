-- name: GetAuthor :one
SELECT  *
FROM    authors a
WHERE   p.id = ?
LIMIT   1;
