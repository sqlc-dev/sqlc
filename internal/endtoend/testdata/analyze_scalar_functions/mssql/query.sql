-- name: ItemTotals :many
SELECT id, dbo.f_add(id, price) AS total, dbo.f_add(@n, DEFAULT) AS d, dbo.label(@s) AS l FROM items;
