-- name: GetProduct :one
SELECT id, name, price, tags
FROM products
WHERE id = $1;

-- name: ListProducts :many
SELECT id, name, price, tags
FROM products
ORDER BY id;

-- name: CreateProduct :one
INSERT INTO products (name, price, tags)
VALUES ($1, $2, $3)
RETURNING id, name, price, tags;
