-- name: PivotDerived :many
SELECT pv.* FROM (SELECT author_id, year, amount FROM sales WHERE author_id > @min_author) s
PIVOT (SUM(amount) FOR year IN ([2023], [2024])) pv;

-- name: PivotTable :many
SELECT * FROM sales PIVOT (COUNT(amount) FOR year IN ([2023], [2024])) AS pv;

-- name: PivotColumns :many
SELECT pv.author_id, pv.[2023] FROM sales PIVOT (MAX(amount) FOR year IN ([2023])) pv
WHERE pv.id = @id;

-- name: Unpivot :many
SELECT u.id, u.metric, u.value FROM stats s
UNPIVOT (value FOR metric IN (q1, q2, q3)) u;

-- name: UnpivotStar :many
SELECT * FROM stats UNPIVOT (value FOR metric IN (q1)) AS u;
