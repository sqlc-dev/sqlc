-- name: AllTypes :many
SELECT * FROM things;

-- name: Casts :one
SELECT
  CAST(amount AS DECIMAL(5,2)) AS a,
  CAST(title AS NVARCHAR(MAX)) AS b,
  CAST(body AS VARCHAR(10)) AS c,
  CONVERT(DATETIME2(3), updated) AS d,
  CAST(amount AS MONEY) AS e,
  CAST(f53 AS FLOAT(24)) AS f,
  CAST(code AS CHAR(3)) AS g
FROM things
WHERE CAST(@a AS DECIMAL(5,2)) > 0
  AND CAST(@b AS NVARCHAR(MAX)) <> N''
  AND CAST(@c AS VARCHAR(10)) <> ''
  AND CONVERT(DATETIME2(3), @d) > '2000-01-01'
  AND CAST(@e AS MONEY) > 0
  AND TRY_CAST(@f AS FLOAT(24)) > 0
  AND CAST(@g AS CHAR(3)) <> '';

-- name: Params :one
SELECT id FROM things
WHERE price = @price AND body = @body AND phone = @phone AND offset_at = @offset_at;
