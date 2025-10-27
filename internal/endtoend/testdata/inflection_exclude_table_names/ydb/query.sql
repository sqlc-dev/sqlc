-- name: DeleteBarByID :one
DELETE FROM bars WHERE id = $id RETURNING id, name;

-- name: DeleteMyDataByID :one
DELETE FROM my_data WHERE id = $id RETURNING id, name;

-- name: DeleteExclusionByID :one
DELETE FROM exclusions WHERE id = $id RETURNING id, name;
