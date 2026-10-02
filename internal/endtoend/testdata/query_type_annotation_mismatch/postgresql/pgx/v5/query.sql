-- name: GetBook :one :type BookWithAuthor
SELECT books.id, books.title, authors.name AS author_name
FROM books
JOIN authors ON authors.id = books.author_id
WHERE books.id = $1;

-- name: ListBooks :many :type BookWithAuthor
SELECT books.id, books.title, books.author_id
FROM books;
