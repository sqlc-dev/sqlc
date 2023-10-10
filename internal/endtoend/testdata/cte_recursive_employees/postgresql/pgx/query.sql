-- name: GetSubordinates :many
WITH RECURSIVE subordinates(name, manager) AS (
    SELECT
        NULL, sqlc.arg(name)::TEXT
    UNION
    SELECT
        s.manager, e.name
    FROM 
        subordinates AS s
    LEFT OUTER JOIN
        employees AS e
    ON
        e.manager = s.manager
    WHERE
        s.manager IS NOT NULL
)
SELECT 
    s.name
FROM
    subordinates AS s
WHERE
    s.name != sqlc.arg(name);