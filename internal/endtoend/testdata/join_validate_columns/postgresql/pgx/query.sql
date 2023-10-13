-- name: ListAuthors :many
SELECT * FROM authors JOIN books ON authors.id = book.author_id1;