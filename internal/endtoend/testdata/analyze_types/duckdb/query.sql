-- name: AllTypes :many
SELECT * FROM things;

-- name: Casts :one
SELECT
  $1::DECIMAL(5,2) AS a,
  $2::INTEGER[] AS b,
  $3::STRUCT(a INTEGER) AS c,
  $4::mood AS d,
  CAST($5 AS VARCHAR(5)) AS e,
  $6::MAP(VARCHAR, INTEGER) AS f,
  $7::INTEGER[3] AS g
FROM things;

-- name: Params :one
SELECT id FROM things
WHERE ints = $1 AND point = $2 AND price = $3 AND m = $4 AND either = $5 AND grid = $6;
