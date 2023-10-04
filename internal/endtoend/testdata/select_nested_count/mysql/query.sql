-- name: GetAuthorsWithBooksCount :many
SELECT *, (
  SELECT COUNT(id) FROM books
  WHERE books.author_id = id
) AS books_count
FROM authors;
