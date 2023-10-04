/* name: ListUsersByID :many */
SELECT first_name, id, last_name FROM users WHERE id < $1;

/* name: ListUserOrders :many */
SELECT
	users.id,
	users.first_name,
	orders.price
FROM
	orders
LEFT JOIN users ON orders.user_id = users.id
WHERE orders.price > sqlc.arg('min_price');

/* name: GetUserByID :one */
SELECT first_name, id, last_name FROM users WHERE id = sqlc.arg('target_id');

/* name: ListUsersByFamily :many */
SELECT first_name, last_name FROM users WHERE age < sqlc.arg('max_age') AND last_name = sqlc.arg('in_family');

/* name: ListUsersWithLimit :many */
SELECT first_name, last_name FROM users LIMIT $1;

/* name: LimitSQLCArg :many */
select first_name, id FROM users LIMIT $1;

/* name: InsertNewUser :exec */
INSERT INTO users (first_name, last_name) VALUES ($1, $2);

/* name: ListUserParenExpr :many */
SELECT * FROM users WHERE (job_status = 'APPLIED' OR job_status = 'PENDING')
AND id > $1
ORDER BY id
LIMIT $2;
