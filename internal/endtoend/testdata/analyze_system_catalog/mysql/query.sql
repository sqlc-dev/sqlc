-- name: ListColumns :many
SELECT table_name, column_name, data_type
FROM information_schema.columns
WHERE table_schema = ?;

-- name: CountTables :one
SELECT count(*) AS total FROM information_schema.tables WHERE table_type = ?;
