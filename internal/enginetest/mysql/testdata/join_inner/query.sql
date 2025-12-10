-- name: GetOrdersWithUser :many
SELECT o.id, o.status, o.total_amount, u.username, u.email
FROM orders o
INNER JOIN users u ON o.user_id = u.id;

-- name: GetOrderItemsWithProduct :many
SELECT oi.quantity, oi.unit_price, p.name AS product_name
FROM order_items oi
INNER JOIN products p ON oi.product_id = p.id;
