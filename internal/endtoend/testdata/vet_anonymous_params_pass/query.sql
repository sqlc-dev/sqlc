-- name: GetByInferredName :one
SELECT id FROM bar
WHERE id = $1;

-- name: GetByArg :one
SELECT id FROM bar
WHERE id = sqlc.arg(target_id);

-- name: GetByAtParam :one
SELECT id FROM bar
WHERE id = @min_id;
