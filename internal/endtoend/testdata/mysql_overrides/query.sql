/* name: GetCount :one */
SELECT id my_id, COUNT(id) id_count FROM users WHERE id > 4;

/* name: GetNameByID :one */
SELECT first_name, last_name FROM users WHERE id = ?;

/* name: GetAll :many */
SELECT * FROM users;

/* name: GetAllUsersOrders :many */
SELECT u.id user_id, u.first_name, o.price, o.id order_id
FROM orders o LEFT JOIN users u ON u.id = o.user_id;

/* name: InsertNewUser :exec */
INSERT INTO users (first_name, last_name) VALUES (?, ?);

/* name: UpdateAllUsers :exec */
update users set first_name = 'Bob';

/* name: UpdateUserAt :exec */
UPDATE users SET first_name = ?, last_name = ? WHERE id > ? AND first_name = ? LIMIT 3;

/* name: InsertUsersFromOrders :exec */
insert into users ( first_name ) select user_id from orders where id = ?;
