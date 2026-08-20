-- name: SyncInventory :exec
MERGE INTO inventory AS t USING inventory_updates AS u ON t.sku = u.sku WHEN MATCHED AND t.quantity <= $1 THEN DELETE WHEN MATCHED THEN UPDATE SET quantity = u.delta WHEN NOT MATCHED THEN INSERT (sku, quantity) VALUES (u.sku, $2);
