/* name: Get :one */
SELECT * FROM books
WHERE book_id = ?;

/* name: Delete :exec */
DELETE FROM books
WHERE book_id = ?;

/* name: ListByTitleYear :many */
SELECT * FROM books
WHERE title = ? AND yr = ?;

/* name: ListByTags :many */
SELECT
  book_id,
  title,
  name,
  isbn,
  tags
FROM books
LEFT JOIN authors ON books.author_id = authors.author_id
WHERE tags = ?;

/* name: Create :execresult */
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

/* name: Update :exec */
UPDATE books
SET title = ?, tags = ?
WHERE book_id = ?;

/* name: UpdateISBN :exec */
UPDATE books
SET title = ?, tags = ?, isbn = ?
WHERE book_id = ?;

/* name: DeleteAuthorBeforeYear :exec */
DELETE FROM books
WHERE yr < ? AND author_id = ?;
-- WHERE yr < sqlc.arg(min_publish_year) AND author_id = ?;
