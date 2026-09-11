-- name: AllTypes :many
SELECT * FROM things;

-- name: StrictTypes :many
SELECT * FROM strict_things;

-- name: Casts :one
SELECT
  CAST(title AS INTEGER) AS a,
  CAST(n AS TEXT) AS b,
  CAST(title AS REAL) AS d,
  weird + 1 AS e,
  title || 'x' AS f
FROM things;

-- name: Params :one
SELECT id FROM things
WHERE title = ? AND price = ? AND untyped = ? AND weird = ? AND big = ? AND n = ?;
