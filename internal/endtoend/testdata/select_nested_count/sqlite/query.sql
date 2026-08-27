-- name: GetAuthorsWithBooksCount :many
SELECT
  *,
  (
    SELECT count(id)
    FROM books
    WHERE books.author_id = id
  ) AS books_count
FROM authors;
