-- name: UpdateSetMultiple :exec
UPDATE foo SET (name, slug) = ($2, $1);
