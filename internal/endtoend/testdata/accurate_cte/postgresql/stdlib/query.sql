-- name: ListExpensiveProducts :many
WITH expensive AS (
    SELECT * FROM products WHERE price > 100
)
SELECT * FROM expensive;

-- name: GetProductStats :one
WITH product_stats AS (
    SELECT COUNT(*) as total, AVG(price) as avg_price FROM products
)
SELECT * FROM product_stats;
