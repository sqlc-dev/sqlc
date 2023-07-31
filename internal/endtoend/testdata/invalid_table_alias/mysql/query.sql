-- https://github.com/sqlc-dev/sqlc/issues/437
CREATE TABLE authors (
  id   INT PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  bio  text
);

-- name: GetAuthor :one
SELECT  *
FROM    authors a
WHERE   p.id = ?
LIMIT   1;
