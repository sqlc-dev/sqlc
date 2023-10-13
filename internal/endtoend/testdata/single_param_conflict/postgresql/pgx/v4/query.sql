-- name: GetAuthorIDByID :one
SELECT  id
FROM    authors
WHERE   id = $1
LIMIT   1;

-- name: GetAuthorByID :one
SELECT  id, name, bio
FROM    authors
WHERE   id = $1
LIMIT   1;

-- name: GetUser :one
SELECT  sub
FROM    users
WHERE   sub = $1
LIMIT   1;

-- https://github.com/sqlc-dev/sqlc/issues/1235

-- name: SetDefaultName :one
UPDATE  authors
SET     name = 'Default Name'
WHERE   id = $1
RETURNING id;
