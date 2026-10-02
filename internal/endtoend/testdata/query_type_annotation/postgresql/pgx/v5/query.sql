-- name: GetBook :one :type BookWithAuthor
SELECT books.id, books.title, authors.name AS author_name
FROM books
JOIN authors ON authors.id = books.author_id
WHERE books.id = $1;

-- name: ListBooksByAuthor :many :type BookWithAuthor
SELECT books.id, books.title, authors.name AS author_name
FROM books
JOIN authors ON authors.id = books.author_id
WHERE books.author_id = $1;

-- name: ListBooksByTitle :many :type BookWithAuthor
SELECT b.id, b.title, a.name AS author_name
FROM books b
JOIN authors a ON a.id = b.author_id
WHERE b.title = $1;

-- name: GetBookAndAuthor :one :type BookAndAuthor
SELECT sqlc.embed(books), sqlc.embed(authors)
FROM books
JOIN authors ON authors.id = books.author_id
WHERE books.id = $1;

-- name: ListBooksAndAuthors :many :type BookAndAuthor
SELECT sqlc.embed(books), sqlc.embed(authors)
FROM books
JOIN authors ON authors.id = books.author_id
ORDER BY books.id;

-- name: ListAuthors :many :type AuthorSummary
SELECT id, name FROM authors;
