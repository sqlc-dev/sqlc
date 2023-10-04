-- name: UpdateSet :exec
UPDATE foo SET name = $2 WHERE slug = $1;
