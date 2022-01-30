CREATE TABLE authors (
  id   bigint    PRIMARY KEY,
  name text      NOT NULL,
  bio  text
);

CREATE TABLE books (
  id bigint PRIMARY KEY,
  author_id bigint NOT NULL
    REFERENCES authors(id),
  title text NOT NULL
);

-- name: GetAuthorsWithBooksCount :many
SELECT *, (
  SELECT COUNT(id) FROM books
  WHERE books.author_id = id
) AS books_count
FROM authors;
