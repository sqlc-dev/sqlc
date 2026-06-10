-- name: GetAuthorName :one
SELECT name FROM authors WHERE id = $1;
