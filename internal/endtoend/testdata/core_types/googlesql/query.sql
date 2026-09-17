-- name: AllTypes :many
SELECT * FROM things;

-- name: Casts :one
SELECT
  CAST(@a AS NUMERIC) AS a,
  CAST(@b AS ARRAY<INT64>) AS b,
  CAST(@d AS STRING) AS d,
  SAFE_CAST(@e AS FLOAT64) AS e
FROM things;

-- name: Params :one
SELECT id FROM things
WHERE s = @s AND n = @n AND smax = @smax AND f32 = @f32;
