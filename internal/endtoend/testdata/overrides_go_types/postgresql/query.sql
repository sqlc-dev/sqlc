-- name: LoadNewStyle :many
SELECT * FROM new_style WHERE id = $1;

-- name: LoadOldStyle :many
SELECT * FROM old_style WHERE id = $1;

