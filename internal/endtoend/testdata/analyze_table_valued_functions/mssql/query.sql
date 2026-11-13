-- name: ApplyBoth :many
SELECT o.id, o.name, o.price, m.id AS mid, m.label FROM dbo.items_over(@p) o CROSS APPLY dbo.ms_items(o.id) m;

-- name: InlineStar :many
SELECT * FROM items_over(@min);

-- name: OuterApply :many
SELECT i.id, m.label FROM items i OUTER APPLY ms_items(i.id) m WHERE i.id = @id;
