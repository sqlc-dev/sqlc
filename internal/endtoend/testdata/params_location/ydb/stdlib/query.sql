-- name: ListUsersByID :many
SELECT first_name, id, last_name FROM users WHERE id < $id;

-- name: ListUserOrders :many
SELECT
	users.id,
	users.first_name,
	orders.price
FROM
	orders
LEFT JOIN users ON orders.user_id = users.id
WHERE orders.price > $min_price;

-- name: GetUserByID :one
SELECT first_name, id, last_name FROM users WHERE id = $target_id;

-- name: ListUsersByFamily :many
SELECT first_name, last_name FROM users WHERE age < $max_age AND last_name = $in_family;

-- name: ListUsersWithLimit :many
SELECT first_name, last_name FROM users LIMIT $limit;

-- name: LimitSQLCArg :many
select first_name, id FROM users LIMIT $limit;

-- name: InsertNewUser :exec
INSERT INTO users (first_name, last_name) VALUES ($first_name, $last_name);

-- name: ListUserParenExpr :many
SELECT * FROM users WHERE (job_status = 'APPLIED' OR job_status = 'PENDING')
AND id > $id
ORDER BY id
LIMIT $limit;












