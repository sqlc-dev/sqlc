-- name: GetBook :one :type
SELECT id, title FROM books WHERE id = $1;

-- name: ListBooks :many :type 9Books
SELECT id, title FROM books;

-- name: ListAuthors :many :type AuthorRow extra
SELECT id, name FROM authors;
