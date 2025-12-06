-- name: GetCastToFloat :many
SELECT id, CAST(amount AS Float32) AS amount FROM data;

-- name: GetCastToInt :many
SELECT id, CAST(quantity AS UInt32) AS quantity FROM data;

-- name: GetCastToDate :many
SELECT id, CAST(created_date AS Date) AS date FROM data;

-- name: GetMultipleCasts :many
SELECT id, CAST(amount AS Float32) AS amount_float, CAST(quantity AS UInt32) AS quantity_int FROM data;
