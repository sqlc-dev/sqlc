-- name: ListAuthors :one
SELECT  *
FROM    authors
WHERE   email = CASE WHEN Cast($email as Text?) = '' THEN NULL ELSE $email END
        OR username = CASE WHEN Cast($username as Text?) = '' THEN NULL ELSE $username END 
LIMIT   1;
