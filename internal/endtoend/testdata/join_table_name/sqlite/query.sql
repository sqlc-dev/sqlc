CREATE TABLE bar (id integer not null);
CREATE TABLE foo (id integer not null, bar integer references bar(id));

-- name: TableName :one
SELECT foo.id
FROM foo
JOIN bar ON foo.bar = bar.id
WHERE bar.id = ? AND foo.id = ?;
