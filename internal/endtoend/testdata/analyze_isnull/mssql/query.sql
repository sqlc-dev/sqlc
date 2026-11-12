-- name: Defaults :many
SELECT id,
    ISNULL(note, raw) AS note_or_raw,
    ISNULL(note, note) AS still_nullable,
    ISNULL(amount, 0) AS amount
FROM readings
WHERE ISNULL(note, raw) = @value;
