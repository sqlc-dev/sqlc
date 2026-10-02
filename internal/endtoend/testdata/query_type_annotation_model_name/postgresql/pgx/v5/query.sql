-- name: ListBookTitles :many :type Book
SELECT books.id, books.title
FROM books;
