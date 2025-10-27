-- name: GetAuthor :one
SELECT * FROM authors
WHERE author_id = $author_id LIMIT 1;

-- name: GetBook :one
SELECT * FROM books
WHERE book_id = $book_id LIMIT 1;

-- name: DeleteBook :exec
DELETE FROM books
WHERE book_id = $book_id;

-- name: BooksByTitleYear :many
SELECT * FROM books
WHERE title = $title AND year = $year;

-- name: BooksByTags :many
SELECT 
  book_id,
  title,
  name,
  isbn,
  tag
FROM books
LEFT JOIN authors ON books.author_id = authors.author_id
WHERE tag IN sqlc.slice(tags);

-- name: CreateAuthor :one
INSERT INTO authors (name) 
VALUES ($name)
RETURNING *;

-- name: CreateBook :one
INSERT INTO books (
    author_id,
    isbn,
    book_type,
    title,
    year,
    available,
    tag
) VALUES (
    $author_id,
    $isbn,
    $book_type,
    $title,
    $year,
    $available,
    $tag
)
RETURNING *;

-- name: UpdateBook :exec
UPDATE books
SET title = $title, tag = $tag
WHERE book_id = $book_id;

-- name: UpdateBookISBN :exec
UPDATE books
SET title = $title, tag = $tag, isbn = $isbn
WHERE book_id = $book_id;


