CREATE TABLE foo (bar BOOLEAN NOT NULL);

-- name: SelectColumnCast :many
SELECT CAST(bar AS UNSIGNED) FROM foo;
