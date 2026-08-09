-- name: ListColumns :many
SELECT table_name, column_name, data_type
FROM information_schema.columns
WHERE table_schema = $1;

-- name: CountRelations :one
SELECT count(*) AS total FROM pg_catalog.pg_class WHERE relkind = $1;
