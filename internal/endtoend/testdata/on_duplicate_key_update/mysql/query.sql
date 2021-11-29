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

-- name: UpsertAuthorNamed :exec
INSERT INTO authors (name, bio)
VALUES (?, sqlc.arg(bio))
ON DUPLICATE KEY
    UPDATE bio = sqlc.arg(bio);
