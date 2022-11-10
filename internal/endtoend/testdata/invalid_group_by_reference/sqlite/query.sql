CREATE TABLE authors (
  id   integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  name text    NOT NULL,
  bio  text,
  UNIQUE(name)
);

-- name: ListAuthors :many
SELECT   *
FROM     authors
GROUP BY invalid_reference;
