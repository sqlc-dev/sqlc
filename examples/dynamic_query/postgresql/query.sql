-- name: GetProducts :many
SELECT * FROM products
WHERE 1=1
sqlc.optional('Category', 'AND category = $1')
sqlc.optional('MinPrice', 'AND price >= $2')
sqlc.optional('IsAvailable', 'AND is_available = $3');

-- name: AddProduct :one
INSERT INTO products (name, category, price, is_available)
VALUES ($1, $2, $3, $4)
RETURNING *;
