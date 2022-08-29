CREATE TABLE bar (
  id serial not null,
  name1 text not null,
  name2 text not null,
  name3 text not null,
  primary key (id));

-- name: DeleteBarByID :execrows
DELETE FROM bar WHERE id = $1;

-- name: DeleteBarByIDAndName :execrows
DELETE FROM bar
WHERE id = $1
AND name1 = $2
AND name2 = $3
AND name3 = $4
;
