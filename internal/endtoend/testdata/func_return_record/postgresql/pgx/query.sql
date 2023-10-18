-- name: GetData :one
SELECT key, value
FROM my_table, jsonb_each(data)
LIMIT 1;
