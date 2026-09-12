-- name: ListAuthorsByMinID :many
SELECT id, name, bio FROM author WHERE id >= ? ORDER BY name;

-- name: GetAuthorByID :one
SELECT id, name, bio FROM author WHERE id = ?;
