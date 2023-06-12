-- Example queries for sqlc
CREATE TABLE authors (
  id   INT
);

-- name: ListAuthors :many
SELECT id FROM authors
ORDER BY adfadsf;