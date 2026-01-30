-- name: GetBetweenPrices :many
SELECT  *
FROM    products
WHERE   price BETWEEN $min_price AND $max_price;

-- name: GetBetweenPricesTable :many
SELECT  *
FROM    products
WHERE   products.price BETWEEN $min_price AND $max_price;

-- name: GetBetweenPricesTableAlias :many
SELECT  *
FROM    products AS p
WHERE   p.price BETWEEN $min_price AND $max_price;

-- name: GetBetweenPricesNamed :many
SELECT  *
FROM    products
WHERE   price BETWEEN sqlc.arg(min_price) AND sqlc.arg(max_price);
