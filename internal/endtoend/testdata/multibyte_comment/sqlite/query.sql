-- Multi-byte UTF-8 in comments must not shift the byte offsets used to slice
-- query text out of the source file (#4523, #4235, #4372).

-- name: ListItems :many
-- an em dash right here — must not truncate the ORDER BY below
SELECT id, name, cap_read
FROM items WHERE cap_read = 'anonymous' ORDER BY name;

-- section — divider between queries

-- name: UpdateItem :exec
UPDATE items SET name = ? WHERE id = ?;

-- name: CreateItem :exec
-- 当新建条目时，需要写入名称和权限
INSERT INTO items (name, cap_read)
VALUES (?, ?);

-- name: DeleteItem :one
DELETE FROM items WHERE id = ? RETURNING id, name, cap_read;
