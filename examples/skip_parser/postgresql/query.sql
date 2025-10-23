-- name: GetProduct :one
SELECT id, name, price, created_at, tags, metadata
FROM products
WHERE id = $1;

-- name: ListProducts :many
SELECT id, name, price, created_at, tags, metadata
FROM products
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CreateProduct :one
INSERT INTO products (name, price, tags, metadata)
VALUES ($1, $2, $3, $4)
RETURNING id, name, price, created_at, tags, metadata;

-- name: UpdateProduct :one
UPDATE products
SET name = $2, price = $3, tags = $4, metadata = $5
WHERE id = $1
RETURNING id, name, price, created_at, tags, metadata;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;

-- name: SearchProductsByTag :many
SELECT id, name, price, created_at, tags, metadata
FROM products
WHERE $1 = ANY(tags)
ORDER BY created_at DESC;

-- name: CountProducts :one
SELECT COUNT(*) FROM products;
