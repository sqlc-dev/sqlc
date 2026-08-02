-- name: MergeReturningOne :one
MERGE INTO inventory AS t
USING inventory_updates AS u
ON t.sku = u.sku AND t.sku = $1
WHEN MATCHED THEN
    UPDATE SET quantity = u.delta
WHEN NOT MATCHED THEN
    INSERT (sku, quantity) VALUES (u.sku, u.delta)
RETURNING t.sku, t.quantity;

-- name: MergeReturningMany :many
MERGE INTO inventory AS t
USING inventory_updates AS u
ON t.sku = u.sku
WHEN MATCHED THEN
    UPDATE SET quantity = u.delta
WHEN NOT MATCHED THEN
    INSERT (sku, quantity) VALUES (u.sku, u.delta)
RETURNING t.sku, t.quantity, t.updated_by;

-- name: MergeReturningStar :many
MERGE INTO product_catalog AS t
USING product_imports AS u
ON t.id = u.label
WHEN MATCHED THEN
    UPDATE SET price = u.id
WHEN NOT MATCHED THEN
    INSERT (id, name, price) VALUES (u.label, u.label, u.id)
RETURNING *;

-- name: MergeReturningAction :many
MERGE INTO inventory AS t
USING inventory_updates AS u
ON t.sku = u.sku
WHEN MATCHED THEN
    UPDATE SET quantity = u.delta
WHEN NOT MATCHED THEN
    INSERT (sku, quantity) VALUES (u.sku, u.delta)
RETURNING merge_action(), t.sku, t.quantity;

-- name: MergeReturningNotMatchedBySource :many
MERGE INTO inventory AS t
USING inventory_updates AS u
ON t.sku = u.sku
WHEN NOT MATCHED BY SOURCE THEN
    DELETE
RETURNING u.sku, u.delta, t.sku, t.quantity;

-- name: MergeReturningSubquerySource :many
MERGE INTO inventory AS t
USING (SELECT sku, delta FROM inventory_updates WHERE delta > 0) AS u
ON t.sku = u.sku
WHEN NOT MATCHED BY SOURCE THEN
    DELETE
RETURNING u.sku, u.delta, t.sku, t.quantity;
