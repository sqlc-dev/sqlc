CREATE TABLE bar (id serial not null, name text not null, primary key (id));

-- name: DeleteBarByID :execrows
DELETE FROM bar WHERE id = $1;

-- name: DeleteBarByIDAndName :execrows
DELETE FROM bar WHERE id = $1 AND name = $2;
