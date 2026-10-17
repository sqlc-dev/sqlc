-- name: GetTimezones :many
SELECT * FROM pg_catalog.pg_timezone_names;

-- name: GetTables :many
SELECT table_name::text FROM information_schema.tables;
