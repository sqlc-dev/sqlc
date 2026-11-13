-- name: Recent :many
SELECT id,
    DATEADD(day, 1, created) AS due,
    DATEDIFF(minute, created, SYSUTCDATETIME()) AS age,
    DATEPART(month, created) AS month_of,
    DATENAME(weekday, created) AS weekday_of
FROM authors
WHERE created > DATEADD(day, @days, SYSUTCDATETIME());

-- name: FromStrings :one
SELECT DATEADD(day, 1, '2026-01-01') AS next_day,
    DATEADD(day, 1, N'2026-01-01') AS next_day_n,
    DATETRUNC(month, '2026-01-15') AS month_start;
