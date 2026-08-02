-- name: MergeWithoutReturning :one
MERGE INTO inventory AS t
USING inventory_updates AS u
ON t.sku = u.sku
WHEN MATCHED THEN
    UPDATE SET quantity = u.delta;
