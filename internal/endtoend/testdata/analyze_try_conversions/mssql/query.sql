-- name: Converted :many
SELECT id,
    TRY_CAST(raw AS INT) AS as_int,
    TRY_CONVERT(DECIMAL(5,1), raw) AS as_decimal,
    TRY_CAST(note AS NVARCHAR(MAX)) AS as_text
FROM readings
WHERE TRY_CAST(@min AS INT) > 0;

-- name: Parsed :many
SELECT id,
    TRY_PARSE(raw AS DECIMAL(10,2)) AS try_parsed,
    TRY_PARSE(note AS DATETIME2 USING 'en-US') AS try_parsed_date
FROM readings
WHERE PARSE(@limit AS INT) > 0;
