-- Example queries for sqlc
CREATE TABLE authors (
  id   BIGSERIAL PRIMARY KEY,
  name TEXT      NOT NULL,
  bio  text
);

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

-- https://github.com/kyleconroy/sqlc/issues/1290
CREATE TABLE users (
  sub UUID PRIMARY KEY
);

-- name: GetUser :one
SELECT  sub
FROM    users
WHERE   sub = $1
LIMIT   1;

-- https://github.com/kyleconroy/sqlc/issues/1235

-- name: SetDefaultName :one
UPDATE  authors
SET     name = "Default Name"
WHERE   id = $1
RETURNING id;
