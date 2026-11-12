-- name: Combined :many
SELECT GREATEST(i, d52) AS greatest_of,
    LEAST(d52, i, si) AS least_of,
    CHOOSE(@pick, i, d102) AS chosen,
    CASE WHEN i > 0 THEN si ELSE d102 END AS case_of,
    IIF(i > 0, v20, n10) AS iif_of,
    COALESCE(m, d102) AS coalesce_of,
    ISNULL(i, d52) AS isnull_int,
    ISNULL(n10, c3) AS isnull_text
FROM t;
