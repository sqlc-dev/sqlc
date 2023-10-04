-- name: GetBetweenPrices :many
SELECT  *
FROM    products
WHERE   price BETWEEN ? AND ?;

-- name: GetBetweenPricesTable :many
SELECT  *
FROM    products
WHERE   products.price BETWEEN ? AND ?;

-- name: GetBetweenPricesTableAlias :many
SELECT  *
FROM    products as p
WHERE   p.price BETWEEN ? AND ?;

-- name: GetBetweenPricesNamed :many
SELECT  *
FROM    products
WHERE   price BETWEEN sqlc.arg(min_price) AND sqlc.arg(max_price);
