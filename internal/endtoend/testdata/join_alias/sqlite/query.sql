CREATE TABLE foo (id integer not null);
CREATE TABLE bar (id integer not null references foo(id), title text);

-- name: AliasJoin :many
SELECT f.id, b.title
FROM foo f
JOIN bar b ON b.id = f.id
WHERE f.id = ?;

-- name: AliasExpand :many
SELECT *
FROM foo f
JOIN bar b ON b.id = f.id
WHERE f.id = ?;
