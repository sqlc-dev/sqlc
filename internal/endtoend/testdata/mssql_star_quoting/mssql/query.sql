-- name: ListOrders :many
SELECT o.* FROM [order] o;

-- name: SalesByYear :many
SELECT pv.* FROM sales PIVOT (SUM(amount) FOR year IN ([2023], [2024])) pv;

-- name: ListOrdersDoubleQuoted :many
SELECT "o".* FROM [order] "o";
