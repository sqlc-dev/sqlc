-- name: ListAuthors :many
SELECT id, name, bio FROM author ORDER BY name;

-- name: GetAuthorByID :one
SELECT id, name, bio FROM author WHERE id = ?;
