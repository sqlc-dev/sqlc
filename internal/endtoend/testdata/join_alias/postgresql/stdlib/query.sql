CREATE TABLE foo (id serial not null);
CREATE TABLE bar (id serial not null references foo(id), title text);

-- name: AliasJoin :many
SELECT f.id, b.title
FROM foo f
JOIN bar b ON b.id = f.id
WHERE f.id = $1;

-- name: AliasExpand :many
SELECT *
FROM foo f
JOIN bar b ON b.id = f.id
WHERE f.id = $1;
