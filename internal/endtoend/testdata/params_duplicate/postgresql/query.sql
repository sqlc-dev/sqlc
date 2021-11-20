CREATE TABLE users (
    id INT PRIMARY KEY,
    first_name varchar(255),
    last_name varchar(255)
);

/* name: SelectUserByID :many */
SELECT first_name from
users where (sqlc.arg(id) = id OR sqlc.arg(id) = 0);

/* name: SelectUserByName :many */
SELECT first_name
FROM users
WHERE first_name = sqlc.arg(name)
   OR last_name = sqlc.arg(name);

/* name: SelectUserQuestion :many */
SELECT first_name from
users where ($1 = id OR  $1 = 0);
