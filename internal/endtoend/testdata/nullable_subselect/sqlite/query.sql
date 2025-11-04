-- name: FirstRowFromFooTable :many
SELECT a, (SELECT a FROM foo limit 1) as "first" FROM foo;

-- name: FirstRowFromEmptyTable :many
SELECT a, (SELECT a FROM empty limit 1) as "first" FROM foo;

-- In SQLite, count() and total() return 0 for empty table.
-- https://www.sqlite.org/lang_aggfunc.html

-- name: CountRowsEmptyTable :many
SELECT a, (SELECT count(a) FROM empty) as "count" FROM foo;

-- name: TotalEmptyTable :many
SELECT a, (SELECT total(a) FROM empty) as "total" FROM foo;
