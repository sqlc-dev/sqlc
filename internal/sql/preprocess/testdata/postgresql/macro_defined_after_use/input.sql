-- name: ListA :many
SELECT * FROM sqlc.macro.Rows() r;

-- name: Rows :macro
SELECT 1 AS n;

-- name: ListB :many
SELECT * FROM sqlc.macro.Rows() r;
