-- name: ListCopyAll :many
SELECT * FROM copy_all;

-- name: GetCopySome :one
SELECT id, label FROM copy_some WHERE id = ?;
