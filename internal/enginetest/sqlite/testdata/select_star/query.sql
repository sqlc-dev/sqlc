-- name: GetAllUsers :many
SELECT * FROM users;

-- name: GetAllProducts :many
SELECT * FROM products;

-- name: GetAllFromSubquery :many
SELECT * FROM (SELECT id, username FROM users) t;
