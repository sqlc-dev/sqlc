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
  tags
FROM books
LEFT JOIN authors ON books.author_id = authors.author_id
WHERE tags = ?;

/* name: CreateAuthor :execresult */
INSERT INTO authors (name) VALUES (?);

/* name: CreateBook :execresult */
INSERT INTO books (
    author_id,
    isbn,
    book_type,
    title,
    yr,
    available,
    tags
) VALUES (
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?
);

/* name: UpdateBook :exec */
UPDATE books
SET title = ?, tags = ?
WHERE book_id = ?;

/* name: UpdateBookISBN :exec */
UPDATE books
SET title = ?, tags = ?, isbn = ?
WHERE book_id = ?;

/* name: DeleteAuthorBeforeYear :exec */
DELETE FROM books
WHERE yr < ? AND author_id = ?;
-- WHERE yr < sqlc.arg(min_publish_year) AND author_id = ?;
