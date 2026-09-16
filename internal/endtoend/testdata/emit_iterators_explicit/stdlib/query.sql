-- name: ListAuthors :many
SELECT id, name FROM author ORDER BY name;

-- name: StreamAuthors :stream
SELECT id, name FROM author ORDER BY name DESC;
