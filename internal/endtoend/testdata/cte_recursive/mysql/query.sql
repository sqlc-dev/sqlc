-- name: CTERecursive :many
WITH RECURSIVE cte AS (
        SELECT b.* FROM bar AS b
        WHERE b.id = ?
    UNION ALL
        SELECT b.*
        FROM bar AS b, cte AS c
        WHERE b.parent_id = c.id
) SELECT * FROM cte;
