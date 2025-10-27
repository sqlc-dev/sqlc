-- name: InsertOrders :exec
INSERT INTO orders (id, name)
SELECT id, CASE WHEN CAST($name_do_update AS Bool) THEN $name ELSE s.name END AS name
FROM orders s;
