-- name: CreateUser :one
INSERT INTO users (username, email, created_at)
VALUES ($1, $2, NOW())
RETURNING *;

-- name: CreateProduct :one
INSERT INTO products (name, price, created_at)
VALUES ($1, $2, NOW())
RETURNING *;
