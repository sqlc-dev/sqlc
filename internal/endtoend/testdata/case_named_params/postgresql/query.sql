-- https://github.com/kyleconroy/sqlc/issues/1195

CREATE TABLE authors (
  id   BIGSERIAL PRIMARY KEY,
  username TEXT NULL,
  email TEXT NULL,
  name TEXT  NOT NULL,
  bio  TEXT
);

-- name: ListAuthors :one
SELECT  *
FROM    authors
WHERE   email = CASE WHEN sqlc.arg(email)::text = '' then NULL else sqlc.arg(email)::text END
        OR username = CASE WHEN sqlc.arg(username)::text = '' then NULL else sqlc.arg(username)::text END 
LIMIT   1;
