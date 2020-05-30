CREATE TABLE users (
    id integer NOT NULL AUTO_INCREMENT PRIMARY KEY,
    first_name varchar(255) NOT NULL,
    last_name varchar(255),
    age integer NOT NULL
) ENGINE=InnoDB;

-- name: GetNameByID :one
SELECT first_name, last_name FROM users WHERE id = ?;
