-- name: AliasJoin :many
SELECT f.id, b.title
FROM foo f
JOIN bar b ON b.id = f.id
WHERE f.id = ?;

-- name: AliasExpand :many
SELECT *
FROM foo f
JOIN bar b ON b.id = f.id
WHERE f.id = ?;
