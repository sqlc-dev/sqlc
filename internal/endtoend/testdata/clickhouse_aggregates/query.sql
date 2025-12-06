-- Basic aggregates
-- name: GetSalesStatistics :many
SELECT
    category,
    COUNT(*) as total_sales,
    SUM(amount) as total_revenue,
    AVG(amount) as avg_amount,
    MIN(amount) as min_amount,
    MAX(amount) as max_amount,
    SUM(quantity) as total_quantity
FROM sales
GROUP BY category
ORDER BY total_revenue DESC;

-- Conditional aggregates
-- name: GetCategoryStats :many
SELECT
    category,
    COUNT(*) as all_sales,
    countIf(amount > 100) as high_value_sales,
    sumIf(amount, quantity > 5) as revenue_bulk_orders,
    avgIf(amount, amount > 50) as avg_high_value
FROM sales
WHERE created_at >= ?
GROUP BY category;

-- HAVING clause
-- name: GetTopCategories :many
SELECT
    category,
    COUNT(*) as sale_count,
    SUM(amount) as total_amount
FROM sales
GROUP BY category
ORDER BY total_amount DESC;

-- Multiple GROUP BY columns
-- name: GetProductCategoryStats :many
SELECT
    product_id,
    category,
    COUNT(*) as count,
    SUM(amount) as total
FROM sales
GROUP BY product_id, category
ORDER BY product_id, total DESC;
