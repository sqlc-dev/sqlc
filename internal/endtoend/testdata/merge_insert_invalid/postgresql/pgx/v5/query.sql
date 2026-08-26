-- name: MergeInsertNoColumnList :exec
MERGE INTO inventory AS t
USING inventory_updates AS u
ON t.sku = u.sku
WHEN NOT MATCHED THEN
    INSERT VALUES (u.sku, $1);
