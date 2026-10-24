-- name: ByID :many
SELECT id, name FROM items WHERE id = sqlc.arg(want_id);

-- name: ByNullableID :many
SELECT id, name FROM items WHERE id = sqlc.narg(want_id);

-- name: ByAt :many
SELECT id, name FROM items WHERE id = @want_id;

-- name: Range :many
SELECT id, name FROM items WHERE id > sqlc.arg(min_id) AND id < sqlc.arg(max_id);
