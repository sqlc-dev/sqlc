CREATE TABLE users (
    id INT PRIMARY KEY,
    name VARCHAR(255)
);

-- name: FindByID :many
SELECT * FROM users WHERE ? = id;

-- name: FindByIDAndName :many
SELECT * FROM users WHERE ? = id AND ? = name;
