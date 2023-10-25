CREATE TABLE users (
    id integer NOT NULL AUTO_INCREMENT PRIMARY KEY,
    first_name varchar(255) NOT NULL,
    last_name varchar(255),
    age integer NOT NULL,
    job_status varchar(10) NOT NULL
);

CREATE TABLE orders (
    id integer NOT NULL AUTO_INCREMENT PRIMARY KEY,
    price DECIMAL(13, 4) NOT NULL,
    user_id integer NOT NULL
);

-- name: SelectUsers :many
SELECT first_name, last_name FROM users WHERE age > sqlc.arg(age);
-- name: SelectUsersDynamic :many
SELECT first_name, last_name FROM users WHERE age > sqlc.arg(age) AND sqlc.dynamic('dynamic');

-- name: SelectUsersDynamic2 :many
SELECT first_name, last_name
FROM users
WHERE sqlc.dynamic('dynamic') AND
    age > sqlc.arg(age) AND
    job_status = sqlc.arg(status) ;
