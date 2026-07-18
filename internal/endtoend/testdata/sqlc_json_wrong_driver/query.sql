-- name: GetOne :one
SELECT sqlc.jsonb_build_object."ItemView"('x', items.x) AS obj FROM items LIMIT 1;
