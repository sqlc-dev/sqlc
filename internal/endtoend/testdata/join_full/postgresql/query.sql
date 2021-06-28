CREATE TABLE foo (id serial not null, bar_id int references bar(id));
CREATE TABLE bar (id serial not null);

-- name: FullJoin :many
SELECT f.id, f.bar_id, b.id
FROM foo f
FULL OUTER JOIN bar b ON b.id = f.bar_id
WHERE f.id = $1;