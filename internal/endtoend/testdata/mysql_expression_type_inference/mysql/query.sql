-- name: FloatDivByConst :many
SELECT (value / 1024) AS scaled_value FROM metrics;

-- name: IntDivByConst :many
SELECT (count / 10) AS avg_count FROM metrics;

-- name: FloatDivByFloat :many
SELECT (value / ratio) AS proportion FROM metrics;

-- name: NotNullFloatDivByConst :many
SELECT (CAST(value AS FLOAT) / 1024) AS scaled FROM metrics;
