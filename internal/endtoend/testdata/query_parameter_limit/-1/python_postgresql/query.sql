-- name: DeleteBarByID :execrows
DELETE FROM bar WHERE id = $1;

-- name: DeleteBarByIDAndName :execrows
DELETE FROM bar WHERE id = $1 AND name = $2;
