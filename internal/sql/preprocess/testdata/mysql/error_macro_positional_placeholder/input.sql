-- name: ByID :macro
SELECT * FROM foo WHERE id = ?;

-- name: Use :many
SELECT * FROM sqlc.macro.ByID() f;
