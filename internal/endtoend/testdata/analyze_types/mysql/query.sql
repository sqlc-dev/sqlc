-- name: AllTypes :many
SELECT * FROM things;

-- name: Casts :one
SELECT
  CAST(count AS UNSIGNED) AS a,
  CAST(price AS DECIMAL(5,2)) AS b,
  CAST(title AS CHAR(10)) AS c,
  CAST(count AS SIGNED) AS d,
  CAST(doc AS JSON) AS e,
  CAST(created AS DATETIME(3)) AS f,
  CAST(price AS DECIMAL) AS g,
  CAST(key16 AS BINARY(8)) AS h
FROM things;

-- name: Params :one
SELECT id FROM things
WHERE count = ? AND price = ? AND kind = ? AND flags = ? AND title = ? AND uprice = ? AND flag = ? AND created = ?;
