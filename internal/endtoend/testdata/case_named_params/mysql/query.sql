-- name: ListAuthors :one
SELECT  *
FROM    authors
WHERE   email = CASE WHEN sqlc.arg(email) = '' then NULL else sqlc.arg(email) END
        OR username = CASE WHEN sqlc.arg(username) = '' then NULL else sqlc.arg(username) END 
LIMIT   1;
