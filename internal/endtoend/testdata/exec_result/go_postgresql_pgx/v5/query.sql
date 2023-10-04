-- name: DeleteBarByID :execresult
DELETE FROM bar WHERE id = $1;
