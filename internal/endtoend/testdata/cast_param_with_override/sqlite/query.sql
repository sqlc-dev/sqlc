-- name: GetData :many
SELECT baz, CAST(max(bar) AS REAL) AS result FROM foo;
