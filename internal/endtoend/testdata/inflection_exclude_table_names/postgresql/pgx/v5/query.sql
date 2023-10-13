-- name: DeleteBarByID :one
DELETE FROM bars WHERE id = $1 RETURNING id, name;

-- name: DeleteMyDataByID :one
DELETE FROM my_data WHERE id = $1 RETURNING id, name;

-- name: DeleteExclusionByID :one
DELETE FROM exclusions WHERE id = $1 RETURNING id, name;
