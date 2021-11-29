CREATE TABLE users (
    id SERIAL NOT NULL,
    first_name varchar(255) NOT NULL,
    last_name varchar(255)
);

/* name: GetUserByID :one */
SELECT first_name, id, last_name FROM users WHERE id = sqlc.arg('target_id');

/* name: ListUsers :many */
SELECT first_name, last_name FROM users;

/* name: InsertNewUser :exec */
INSERT INTO users (first_name, last_name) VALUES (?, ?);

/* name: InsertNewUserWithResult :execresult */
INSERT INTO users (first_name, last_name) VALUES (?, ?);

/* name: DeleteUsersByName :execrows */
DELETE FROM users WHERE first_name = ? AND last_name = ?;
