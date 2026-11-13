-- name: Defaults :many
SELECT id,
    ISNULL(note, raw) AS note_or_raw,
    ISNULL(note, note) AS still_nullable,
    ISNULL(amount, 0) AS amount
FROM readings
WHERE ISNULL(note, raw) = @value;

-- name: FirstArgumentTypes :many
SELECT ISNULL(raw, CAST(note AS NVARCHAR(40))) AS first_text,
    COALESCE(raw, CAST(note AS NVARCHAR(40))) AS combined_text,
    ISNULL(NULL, raw) AS literal_null,
    ISNULL(@fallback, raw) AS placeholder_first,
    ISNULL(note, @replacement) AS placeholder_second
FROM readings;
