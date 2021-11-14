-- https://github.com/kyleconroy/sqlc/issues/921
CREATE TABLE authors (
  id   BIGSERIAL PRIMARY KEY,
  name text      NOT NULL UNIQUE,
  bio  text
);

-- name: UpsertAuthor :exec
INSERT INTO authors (name, bio)
VALUES ($1, $2)
ON CONFLICT (name) DO UPDATE
SET bio = $2;
