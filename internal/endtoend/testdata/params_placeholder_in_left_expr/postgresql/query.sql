CREATE TABLE users (
    id INT PRIMARY KEY,
    name VARCHAR(255)
);

-- name: FindByID :many
SELECT * FROM users WHERE $1 = id;

-- name: FindByIDAndName :many
SELECT * FROM users WHERE $1 = id AND $1 = name;
