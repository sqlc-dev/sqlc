CREATE TABLE users (
    first_name varchar(255),
    last_name  varchar(255),
    age        smallint
);

-- name: GetAll :many
SELECT * FROM users;
