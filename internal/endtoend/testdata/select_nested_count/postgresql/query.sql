CREATE TABLE authors (
  id   BIGSERIAL PRIMARY KEY,
  name text      NOT NULL,
  bio  text
);

CREATE TABLE books (
  id BIGSERIAL PRIMARY KEY,
  author_id BIGSERIAL NOT NULL
  	REFERENCES authors(id),
  title text NOT NULL
);

-- name: GetAuthorsWithBooksCount :many
SELECT *, (
  SELECT COUNT(id) FROM books
  WHERE books.author_id = id
) AS books_count
FROM authors;
