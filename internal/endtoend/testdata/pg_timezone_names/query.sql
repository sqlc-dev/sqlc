-- name: GetTimezones :many
SELECT * from pg_catalog.pg_timezone_names;

-- name: GetTables :many
SELECT table_name::text from information_schema.tables;

-- name: GetColumns :many
SELECT table_name::text, column_name::text from information_schema.columns;
