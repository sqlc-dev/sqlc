-- name: UpdateAuthor :exec
UPDATE authors
SET name = $1
WHERE id = $1;
