-- https://github.com/kyleconroy/sqlc/issues/1195

CREATE TABLE authors (
  id   BIGINT PRIMARY KEY,
  username TEXT NULL,
  email TEXT NULL,
  name TEXT  NOT NULL,
  bio  TEXT,
  UNIQUE KEY idx_username (username),
  UNIQUE KEY ids_email (email)
);

-- name: ListAuthors :one
SELECT  *
FROM    authors
WHERE   email = CASE WHEN sqlc.arg(email) = '' then NULL else sqlc.arg(email) END
        OR username = CASE WHEN sqlc.arg(username) = '' then NULL else sqlc.arg(username) END 
LIMIT   1;
