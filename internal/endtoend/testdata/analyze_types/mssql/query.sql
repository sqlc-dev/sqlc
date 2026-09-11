-- name: AllTypes :many
SELECT * FROM things;

-- name: Casts :one
SELECT
  CAST(@a AS DECIMAL(5,2)) AS a,
  CAST(@b AS NVARCHAR(MAX)) AS b,
  CAST(@c AS VARCHAR(10)) AS c,
  CONVERT(DATETIME2(3), @d) AS d,
  CAST(@e AS dbo.PhoneNumber) AS e,
  TRY_CAST(@f AS FLOAT(24)) AS f,
  CAST(@g AS dbo.Code) AS g
FROM things;

-- name: Params :one
SELECT id FROM things
WHERE price = @price AND body = @body AND phone = @phone AND vec = @vec AND offset_at = @offset_at;
