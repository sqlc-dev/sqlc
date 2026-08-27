-- name: AliasJoin :many
SELECT f.id, b.title
FROM foo AS f
JOIN bar AS b ON b.id = f.id
WHERE f.id = ?;

-- name: AliasExpand :many
SELECT *
FROM foo AS f
JOIN bar AS b ON b.id = f.id
WHERE f.id = ?;
