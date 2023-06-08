CREATE TABLE bar (id serial not null);

-- name: DeleteBarByID :execrows
DELETE FROM bar WHERE id = $1;
