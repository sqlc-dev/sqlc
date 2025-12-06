-- name: GetProductsWithoutDescription :many
SELECT id, name FROM products WHERE description IS NULL;

-- name: GetProductsWithDescription :many
SELECT id, name, description FROM products WHERE description IS NOT NULL;

-- name: GetProductsWithDiscount :many
SELECT id, name, coalesce(discount, 0) AS discount FROM products;

-- name: GetProductsWithDefault :many
SELECT id, name, ifNull(category, 'Uncategorized') AS category FROM products;
