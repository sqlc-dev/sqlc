-- Example queries for sqlc
CREATE TABLE authors (
  id   BIGSERIAL PRIMARY KEY,
  name text      NOT NULL,
  bio  text,
  tags string[]
);

-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = $1 LIMIT 1;