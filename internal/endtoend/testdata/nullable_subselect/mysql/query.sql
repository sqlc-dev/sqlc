-- name: FirstRowFromFooTable :many
SELECT a, (SELECT a FROM foo limit 1) as "first" FROM foo;

-- name: FirstRowFromEmptyTable :many
SELECT a, (SELECT a FROM empty limit 1) as "first" FROM foo;

-- In MySQL, only count() returns 0 for empty table.
-- https://dev.mysql.com/doc/refman/8.0/en/aggregate-functions.html
-- name: CountRowsEmptyTable :many
SELECT a, (SELECT count(a) FROM empty) as "count" FROM foo;
