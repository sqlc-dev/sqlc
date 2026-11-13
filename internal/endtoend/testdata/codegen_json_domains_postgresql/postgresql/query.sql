-- name: ListBooks :many
SELECT id, title, price, ref FROM books WHERE price > $1;
