-- name: CreatedIn :many
SELECT id,
    created AT TIME ZONE 'UTC' AS created_utc,
    GETDATE() AT TIME ZONE 'UTC' AS now_utc
FROM authors
WHERE created AT TIME ZONE @tz > SYSDATETIMEOFFSET();
