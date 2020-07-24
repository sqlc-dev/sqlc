CREATE TABLE foo (a text, b text);

-- name: StarExpansion :many
SELECT *, *, foo.* FROM foo;
