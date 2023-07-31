CREATE TABLE foo (bar TEXT NOT NULL);

-- name: SelectColumnCast :many
SELECT CAST(bar AS BLOB) FROM foo;
