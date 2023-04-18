-- https://github.com/kyleconroy/sqlc/issues/1198
CREATE TABLE authors (
  id   INT PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  bio  text
);

-- name: SetAuthor :exec
UPDATE  authors
SET     name = ?
WHERE   id = ?
