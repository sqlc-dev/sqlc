-- name: JoinWhereClause :many
SELECT foo.*
FROM foo
JOIN bar ON bar.id = barid
WHERE owner = $owner;

-- name: JoinParamWhereClause :many
SELECT f.*
FROM foo AS f
JOIN bar AS b ON b.id = f.barid
WHERE b.id = $bar_id AND b.owner = $owner;

-- name: JoinNoConstraints :many
SELECT foo.*
FROM foo
CROSS JOIN bar
WHERE bar.id = $bar_id AND owner = $owner;
