-- name: ListAuthors :one
SELECT  *
FROM    authors
WHERE   email = CASE WHEN sqlc.arg(email)::text = '' then NULL else sqlc.arg(email)::text END
        OR username = CASE WHEN sqlc.arg(username)::text = '' then NULL else sqlc.arg(username)::text END 
LIMIT   1;
