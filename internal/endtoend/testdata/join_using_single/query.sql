-- name: GetOrdersWithShipments :many
SELECT *
FROM orders
LEFT JOIN shipments USING (order_id);
