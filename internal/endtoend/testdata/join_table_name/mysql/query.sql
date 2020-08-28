CREATE TABLE bar (id serial not null);
CREATE TABLE foo (id serial not null, bar serial references bar(id));

-- name: TableName :one
SELECT foo.id
FROM foo
JOIN bar ON foo.bar = bar.id
WHERE bar.id = ? AND foo.id = ?;
