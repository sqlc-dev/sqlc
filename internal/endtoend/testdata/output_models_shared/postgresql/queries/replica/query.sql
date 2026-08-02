-- name: ListAuthors :many
SELECT * FROM authors ORDER BY name;

-- name: ListAuthorsByStatus :many
SELECT * FROM authors WHERE status = $1 ORDER BY name;
