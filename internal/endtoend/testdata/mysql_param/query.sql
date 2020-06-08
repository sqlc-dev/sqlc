/* name: ListUsersByID :many */
SELECT first_name, id, last_name FROM users WHERE id < ?;

/* name: ListUserOrders :many */
SELECT
	users.id,
	users.first_name,
	orders.price
FROM
	orders
LEFT JOIN users ON orders.user_id = users.id
WHERE orders.price > :minPrice;

/* name: GetUserByID :one */
SELECT first_name, id, last_name FROM users WHERE id = :targetID;

/* name: ListUsersByFamily :many */
SELECT first_name, last_name FROM users WHERE age < :maxAge AND last_name = :inFamily;

/* name: ListUsersWithLimit :many */
SELECT first_name, last_name FROM users LIMIT ?;

/* name: LimitSQLCArg :many */
select first_name, id FROM users LIMIT sqlc.arg(UsersLimit);

/* name: InsertNewUser :exec */
INSERT INTO users (first_name, last_name) VALUES (?, sqlc.arg(user_last_name));

/* name: ListUserParenExpr :many */
SELECT * FROM users WHERE (job_status = 'APPLIED' OR job_status = 'PENDING')
AND id > :lastID
ORDER BY id
LIMIT :usersLimit;
