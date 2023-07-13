CREATE TABLE authors (
  name text      NOT NULL,
  bio  text
);

CREATE TABLE people (
  first_name text      NOT NULL
);

-- name: ListAuthorsUnion :many
SELECT name as foo FROM authors
UNION
SELECT first_name as foo FROM people
ORDER BY foo;
