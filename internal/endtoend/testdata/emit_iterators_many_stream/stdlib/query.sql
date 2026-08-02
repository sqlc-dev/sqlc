-- name: ListAuthors :many:stream
SELECT id, name FROM author ORDER BY name;

-- name: ListAllAuthors :many
SELECT id, name FROM author ORDER BY name;
