-- name: Recent :many
SELECT id,
    DATEADD(day, @days, created) AS due,
    DATEDIFF(minute, created, SYSUTCDATETIME()) AS age,
    DATEPART(month, created) AS month_of,
    DATENAME(weekday, created) AS weekday_of
FROM authors
WHERE created > DATEADD(hour, -1, SYSUTCDATETIME());
