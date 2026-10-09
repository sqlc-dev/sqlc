-- name: AllTypes :many
SELECT * FROM things;

-- name: StarColumns :many
SELECT id, name, tag, amount, tags, labels, matrix, kind, created, updated, price, status, attrs, pos, geo, scores, ip, uid, fixed, flag, small, plain, either, n.a, n.b, total, whole FROM things;

-- name: Casts :one
SELECT
  CAST(id AS String) AS a,
  toDecimal64(id, 4) AS b,
  CAST(name AS Nullable(String)) AS c,
  toDateTime64(created, 3) AS d,
  CAST(tags AS Array(String)) AS e,
  CAST(tag AS Nullable(String)) AS f
FROM things;

-- name: Placeholders :many
SELECT id FROM things
WHERE id = {p1:UInt64} AND name = {p2:String} AND amount > {p3:Float64} AND tag = {p4:Nullable(String)} AND price = {p5:Decimal(10, 2)};
