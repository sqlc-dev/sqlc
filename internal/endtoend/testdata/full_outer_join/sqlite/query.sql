-- Example queries for sqlc
CREATE TABLE authors (
  id   INTEGER PRIMARY KEY,
  name text      NOT NULL
);

CREATE TABLE books (
  id INTEGER PRIMARY KEY,
  title text NOT NULL
);

-- name: GetAuthor :one
SELECT * FROM authors AS a
FULL OUTER JOIN books AS b
 ON a.id = b.id
WHERE a.id = ? LIMIT 1;