-- Example queries for sqlc
CREATE TABLE authors (
  id   BIGINT PRIMARY KEY,
  name TEXT      NOT NULL,
  bio  text
);

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

-- https://github.com/kyleconroy/sqlc/issues/1290
CREATE TABLE users (
  sub TEXT PRIMARY KEY
);

-- name: GetUser :one
SELECT  sub
FROM    users
WHERE   sub = ?
LIMIT   1;
