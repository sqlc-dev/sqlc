-- name: GetAuthor :one
SELECT * FROM authors WHERE id = @id;
