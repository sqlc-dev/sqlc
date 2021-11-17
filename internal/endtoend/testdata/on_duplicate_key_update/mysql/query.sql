-- https://github.com/kyleconroy/sqlc/issues/921
CREATE TABLE authors (
  id   BIGINT  NOT NULL AUTO_INCREMENT PRIMARY KEY,
  name text    NOT NULL,
  bio  text,
  UNIQUE(name)
);

-- name: UpsertAuthor :exec
INSERT INTO authors (name, bio)
VALUES (?, ?)
ON DUPLICATE KEY
    UPDATE bio = ?;
