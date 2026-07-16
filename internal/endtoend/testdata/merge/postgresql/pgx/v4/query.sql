-- name: SyncInventory :exec
MERGE INTO inventory AS t
USING inventory_updates AS u
ON t.sku = u.sku
WHEN MATCHED AND t.quantity + u.delta <= $1 THEN
    DELETE
WHEN MATCHED THEN
    UPDATE SET quantity = t.quantity + u.delta, updated_by = $2
WHEN NOT MATCHED THEN
    INSERT (sku, quantity, updated_by) VALUES (u.sku, u.delta, $2);

-- name: MergeFromSubquery :exec
MERGE INTO inventory AS t
USING (SELECT sku, delta FROM inventory_updates WHERE delta > $1) AS u
ON t.sku = u.sku AND t.quantity > $2
WHEN MATCHED THEN
    UPDATE SET quantity = $3
WHEN NOT MATCHED THEN
    INSERT (sku, quantity) VALUES (u.sku, $3);

-- name: MergeDelete :exec
MERGE INTO inventory AS t
USING inventory_updates AS u
ON t.sku = u.sku
WHEN MATCHED THEN
    DELETE;

-- name: MergeNamedParams :exec
MERGE INTO inventory AS t
USING inventory_updates AS u
ON t.sku = u.sku
WHEN MATCHED THEN
    UPDATE SET updated_by = @updated_by
WHEN NOT MATCHED THEN
    INSERT (sku, quantity, updated_by) VALUES (u.sku, sqlc.arg(start_qty), @updated_by);

-- name: MergeMultiSet :exec
MERGE INTO inventory AS t
USING inventory_updates AS u
ON t.sku = u.sku
WHEN MATCHED THEN
    UPDATE SET (quantity, updated_by) = ($1, $2);
