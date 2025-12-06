-- name: ListOrdersByTotal :many
SELECT id, customer_name, total FROM orders ORDER BY total DESC;

-- name: ListOrdersByCustomerAndDate :many
SELECT id, customer_name, created_at FROM orders ORDER BY customer_name, created_at DESC;

-- name: ListOrdersAscending :many
SELECT id, total FROM orders ORDER BY total ASC;
