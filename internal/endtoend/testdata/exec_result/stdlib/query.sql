CREATE TABLE bar (id serial not null);

-- name: DeleteBarByID :execresult
DELETE FROM bar WHERE id = $1;
