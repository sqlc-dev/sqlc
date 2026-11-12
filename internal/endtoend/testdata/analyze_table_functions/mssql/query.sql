-- name: Elements :many
SELECT d.id, j.[key], j.value, j.type
FROM documents d
CROSS APPLY OPENJSON(d.body) j
WHERE d.id = @id;

-- name: Items :many
SELECT d.id, i.name, i.qty, i.extra
FROM documents d
CROSS APPLY OPENJSON(d.body, '$.items')
    WITH (name NVARCHAR(50) '$.name', qty INT '$.qty', extra NVARCHAR(MAX) '$.extra' AS JSON) i;

-- name: Words :many
SELECT d.id, w.value, w.ordinal
FROM documents d
CROSS APPLY STRING_SPLIT(d.title, N' ', 1) w;

-- name: Series :many
SELECT value FROM GENERATE_SERIES(1, 10);
