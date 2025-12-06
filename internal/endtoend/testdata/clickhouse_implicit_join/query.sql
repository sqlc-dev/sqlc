-- name: GetCustomerOrders :many
SELECT c.name, o.amount FROM customers c, orders o WHERE c.id = o.customer_id;

-- name: GetCustomerOrdersByCity :many
SELECT c.name, c.city, o.amount FROM customers c, orders o WHERE c.id = o.customer_id AND c.city = ?;
