CREATE TABLE foo (a text, b text);

-- name: SelectUnion :many
SELECT * FROM foo
UNION
SELECT * FROM foo;
