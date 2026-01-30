-- name: UpdateSet :exec
UPDATE foo SET name = $name WHERE slug = $slug;
