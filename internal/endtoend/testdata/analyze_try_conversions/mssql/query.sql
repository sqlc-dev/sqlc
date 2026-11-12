-- name: Converted :many
SELECT id,
    TRY_CAST(raw AS INT) AS as_int,
    TRY_CONVERT(DECIMAL(5,1), raw) AS as_decimal,
    TRY_CAST(note AS NVARCHAR(MAX)) AS as_text
FROM readings
WHERE TRY_CAST(@min AS INT) > 0;
