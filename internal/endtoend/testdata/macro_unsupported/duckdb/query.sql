-- DuckDB queries are not preprocessed, so there is no sqlc.macro to call a
-- macro with.
-- name: AllFoo :macro
SELECT * FROM foo;

-- name: ListFoo :many
SELECT * FROM foo;
