-- name: GetAuthorIDByID :one
SELECT  id
FROM    authors
WHERE   id = ?
LIMIT   1;

-- name: GetAuthorByID :one
SELECT  id, name, bio
FROM    authors
WHERE   id = ?
LIMIT   1;

-- name: GetUser :one
SELECT  sub
FROM    users
WHERE   sub = ?
LIMIT   1;
