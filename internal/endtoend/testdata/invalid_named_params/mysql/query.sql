-- name: ListAuthors :one
SELECT  *
FROM    authors
WHERE   id = sqlc.arg(my_named_param)
        OR bio = sqlc.arg(my_named_param)
LIMIT   1;
