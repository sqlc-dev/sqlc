-- name: Negated :many
SELECT -CAST(i AS DECIMAL) AS nd, d * CAST(i AS DECIMAL(12,3)) AS product FROM p2;
