-- name: FirstRowFromFooTable :many
SELECT a, (SELECT a FROM foo limit 1) as "first" FROM foo;

-- name: FirstRowFromEmptyTable :many
SELECT a, (SELECT a FROM empty limit 1) as "first" FROM foo;

-- In PostgreSQL, only count() returns 0 for empty table.
-- https://www.postgresql.org/docs/15/functions-aggregate.html
-- name: CountRowsEmptyTable :many
SELECT a, (SELECT count(a) FROM empty) as "count" FROM foo;
