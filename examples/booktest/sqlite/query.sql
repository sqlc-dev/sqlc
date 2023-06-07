/* name: GetAuthor :one */
SELECT * FROM authors
WHERE author_id = ?;

/* name: GetBook :one */
SELECT * FROM books
WHERE book_id = ?;

/* name: DeleteBook :exec */
DELETE FROM books
WHERE book_id = ?;

/* name: BooksByTitleYear :many */
SELECT * FROM books
WHERE title = ? AND yr = ?;

/* name: BooksByTags :many */
SELECT
  book_id,
  title,
  name,
  isbn,
  tag
FROM books
LEFT JOIN authors ON books.author_id = authors.author_id
WHERE tag IN (sqlc.slice(tags));

/* name: CreateAuthor :one */
INSERT INTO authors (name) VALUES (?)
RETURNING *;

/* name: CreateBook :one */
INSERT INTO books (
    author_id,
    isbn,
    book_type,
    title,
    yr,
    available,
    tag
) VALUES (
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?
)
RETURNING *;

/* name: UpdateBook :exec */
UPDATE books
SET title = ?1, tag = ?2
WHERE book_id = ?3;

/* name: UpdateBookISBN :exec */
UPDATE books
SET title = ?1, tag = ?2, isbn = ?4
WHERE book_id = ?3;

/* name: DeleteAuthorBeforeYear :exec */
DELETE FROM books
WHERE yr < ? AND author_id = ?;
