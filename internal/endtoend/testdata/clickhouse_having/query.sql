-- name: GetRegionSalesAboveThreshold :many
SELECT region, sum(amount) AS total_sales FROM sales GROUP BY region HAVING sum(amount) > ?;

-- name: GetRegionSalesCount :many
SELECT region, count(id) AS transaction_count FROM sales GROUP BY region HAVING count(id) > ?;

-- name: GetAverageSalesByRegion :many
SELECT region, avg(amount) AS avg_sale FROM sales GROUP BY region HAVING avg(amount) > ?;
