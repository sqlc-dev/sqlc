-- name: JoinBar :one
SELECT f.id, info
FROM foo f
LEFT JOIN bar b ON b.foo_id = f.id;

-- name: JoinBarAlias :one
SELECT f.id, b.info
FROM foo f
LEFT JOIN bar b ON b.foo_id = f.id;
