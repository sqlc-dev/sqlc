-- name: CreateUser :one
INSERT INTO users (username, email)
VALUES (?, ?)
RETURNING *;

-- name: CreateProduct :one
INSERT INTO products (name, price)
VALUES (?, ?)
RETURNING *;
