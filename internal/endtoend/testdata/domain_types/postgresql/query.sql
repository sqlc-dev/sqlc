-- name: ListBooks :many
SELECT id, title, price, ref FROM books WHERE price > $1;

-- name: ListShelves :many
SELECT * FROM shelves;

-- name: BooksByRef :many
SELECT b.id, b.title, s.label FROM books b LEFT JOIN shelves s ON s.label = b.title WHERE b.ref = $1;
