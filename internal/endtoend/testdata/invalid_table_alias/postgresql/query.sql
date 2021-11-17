-- https://github.com/kyleconroy/sqlc/issues/437
CREATE TABLE authors (
  id   BIGSERIAL PRIMARY KEY,
  name text      NOT NULL,
  bio  text
);

-- name: GetAuthor :one
SELECT  *
FROM    authors a
WHERE   p.id = $1
LIMIT   1;
