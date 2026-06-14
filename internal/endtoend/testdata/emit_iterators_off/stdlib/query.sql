-- name: ListAuthors :many
SELECT id, name, bio FROM author ORDER BY name;
