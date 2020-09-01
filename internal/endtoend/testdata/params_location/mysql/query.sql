CREATE TABLE users (
    id integer NOT NULL AUTO_INCREMENT PRIMARY KEY,
    first_name varchar(255) NOT NULL,
    last_name varchar(255),
    age integer NOT NULL,
    job_status varchar(10) NOT NULL
) ENGINE=InnoDB;

CREATE TABLE orders (
    id integer NOT NULL AUTO_INCREMENT PRIMARY KEY,
    price DECIMAL(13, 4) NOT NULL,
    user_id integer NOT NULL
) ENGINE=InnoDB;

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
WHERE orders.price > sqlc.arg('min_price');

/* name: GetUserByID :one */
SELECT first_name, id, last_name FROM users WHERE id = sqlc.arg('target_id');

/* name: ListUsersByFamily :many */
SELECT first_name, last_name FROM users WHERE age < sqlc.arg('max_age') AND last_name = sqlc.arg('in_family');

/* name: ListUsersWithLimit :many */
SELECT first_name, last_name FROM users LIMIT ?;

/* name: LimitSQLCArg :many */
select first_name, id FROM users LIMIT ?;

/* name: InsertNewUser :exec */
INSERT INTO users (first_name, last_name) VALUES (?, ?);

/* name: ListUserParenExpr :many */
SELECT * FROM users WHERE (job_status = 'APPLIED' OR job_status = 'PENDING')
AND id > ?
ORDER BY id
LIMIT ?;
