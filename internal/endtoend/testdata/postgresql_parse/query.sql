-- name: GetCount :one
SELECT id my_id, COUNT(id) id_count FROM users WHERE id > 4;

-- name: GetNameByID :one
SELECT first_name, last_name FROM users WHERE id = $1;

-- name: GetAll :many
SELECT * FROM users;

-- name: GetAllUsersOrders :many
SELECT u.id user_id, u.first_name, o.price, o.id order_id
FROM orders o
LEFT JOIN users u ON u.id = o.user_id;

-- name: InsertNewUser :exec
INSERT INTO users (first_name, last_name) VALUES ($1, $2);

-- name: UpdateAllUsers :exec
update users set first_name = 'Bob';

-- name: UpdateUserAt :exec
UPDATE users SET first_name = $1, last_name = $2 WHERE id > $3 AND first_name = $1;

-- name: InsertUsersFromOrders :exec
insert into users ( first_name ) select user_id from orders where orders.id = $1;
