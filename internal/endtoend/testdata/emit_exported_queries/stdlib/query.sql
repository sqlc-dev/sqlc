-- name: UpdateBarID :exec
UPDATE bar SET id = $1 WHERE id = $2;
