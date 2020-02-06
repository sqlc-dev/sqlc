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

/* name: CreateAuthor :exec */
INSERT INTO authors (name) VALUES (?);

/* name: CreateBook :exec */
INSERT INTO books (
    author_id,
    isbn,
    booktype,
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
    ?
);

/* name: UpdateBook :exec */
UPDATE books
SET title = ?, tags = ?
WHERE book_id = ?;

/* name: UpdateBookISBN :exec */
UPDATE books
SET title = ?, tags = :book_tags, isbn = ?
WHERE book_id = ?;

/* name: DeleteAuthorBeforeYear :exec */
DELETE FROM books
WHERE yr < sqlc.arg(min_publish_year) AND author_id = ?;