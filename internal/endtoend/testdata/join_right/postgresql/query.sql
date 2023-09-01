CREATE TABLE bar (id serial not null unique);
CREATE TABLE foo (id serial not null, bar_id int references bar(id));

-- name: RightJoin :many
SELECT f.id, f.bar_id, b.id
FROM foo f
RIGHT JOIN bar b ON b.id = f.bar_id
WHERE f.id = $1;
