-- name: GetAuthor :one
SELECT * FROM authors
WHERE author_id = $1;

-- name: GetBook :one
SELECT * FROM books
WHERE book_id = $1;

-- name: DeleteBook :exec
DELETE FROM books
WHERE book_id = $1;

-- name: BooksByTitleYear :many
SELECT * FROM books
WHERE title = $1 AND year = $2;

-- name: BooksByTags :many
SELECT 
  book_id,
  title,
  name,
  isbn,
  tags
FROM books
LEFT JOIN authors ON books.author_id = authors.author_id
WHERE tags && $1::varchar[];

-- name: CreateAuthor :one
INSERT INTO authors (name) VALUES ($1)
RETURNING *;

-- name: CreateBook :one
INSERT INTO books (
    author_id,
    isbn,
    book_type,
    title,
    year,
    available,
    tags
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
RETURNING *;

-- name: UpdateBook :exec
UPDATE books
SET title = $1, tags = $2
WHERE book_id = $3;

-- name: UpdateBookISBN :exec
UPDATE books
SET title = $1, tags = $2, isbn = $4
WHERE book_id = $3;
