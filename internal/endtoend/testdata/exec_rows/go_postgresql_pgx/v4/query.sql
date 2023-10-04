-- name: DeleteBarByID :execrows
DELETE FROM bar WHERE id = $1;
