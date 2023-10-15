CREATE TABLE users (
    id int PRIMARY KEY,
    first_name text NOT NULL,
    last_name text,
    age int NOT NULL,
    job_status text NOT NULL
);

CREATE TABLE orders (
    id int PRIMARY KEY,
    price numeric NOT NULL,
    user_id int NOT NULL
);

-- name: SelectUsers :many
SELECT first_name, last_name FROM users WHERE age > sqlc.arg(age);

-- name: SelectUsersDynamic :many
SELECT first_name, last_name FROM users WHERE age > sqlc.arg(age) AND sqlc.dynamic('dynamic');

-- name: SelectUsersDynamic2 :many
SELECT first_name, last_name
FROM users
WHERE age > sqlc.arg(age) AND
    job_status = sqlc.arg(status) AND
    sqlc.dynamic('dynamic');

-- name: SelectUsersDynamicMulti :many
SELECT first_name, last_name
FROM users
WHERE age > sqlc.arg(age) AND
    job_status = sqlc.arg(status) AND
    sqlc.dynamic('dynamic')
ORDER BY sqlc.dynamic('order');
