CREATE TABLE foo (id serial not null, bar_id int references bar(id));
CREATE TABLE bar (id serial not null);

-- name: FullJoin :many
SELECT *
FROM foo f
LEFT JOIN bar b ON b.id = f.bar_id
UNION ALL
SELECT *
FROM bar b
RIGHT JOIN foo f ON b.id = f.bar_id
WHERE f.id = ?;
