-- name: ListProductsPage :many
SELECT id, name, price FROM products ORDER BY id LIMIT ? OFFSET ?;

-- name: GetFirstNProducts :many
SELECT id, name, price FROM products LIMIT ?;

-- name: GetProductsWithOffset :many
SELECT id, name, price FROM products ORDER BY id LIMIT 10 OFFSET ?;
