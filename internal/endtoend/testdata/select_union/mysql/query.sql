CREATE TABLE foo (a text, b text);

-- name: SelectUnion :many
SELECT * FROM foo
UNION
SELECT * FROM foo;

-- name: SelectExcept :many
SELECT * FROM foo
EXCEPT
SELECT * FROM foo;

-- name: SelectIntersect :many
SELECT * FROM foo
INTERSECT
SELECT * FROM foo;
