-- name: AllTypes :many
SELECT * FROM things;

-- name: Casts :one
SELECT
  CAST(@a AS NUMERIC(5,2)) AS a,
  CAST(@b AS ARRAY<INT64>) AS b,
  CAST(@c AS STRUCT<x INT64>) AS c,
  CAST(@d AS STRING(5)) AS d,
  SAFE_CAST(@e AS BIGNUMERIC) AS e,
  [1, 2] AS f
FROM things;

-- name: Params :one
SELECT id FROM things
WHERE s = @s AND ai = @ai AND n = @n AND st = @st AND smax = @smax;
