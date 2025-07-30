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

-- name: SubqueryAlias :many
SELECT * FROM (SELECT 1 AS n) AS x WHERE x.n <= ?;

-- name: ColumnAlias :many
SELECT * FROM (SELECT 1 AS n) WHERE n <= ?;

-- name: ColumnAndQueryAlias :many
SELECT * FROM (SELECT 1 AS n) AS x WHERE n <= ?;
