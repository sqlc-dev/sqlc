-- name: JoinWhereClause :many
SELECT foo.*
FROM foo
JOIN bar ON bar.id = barid
WHERE owner = $owner;

-- name: JoinParamWhereClause :many
SELECT foo.*
FROM foo
JOIN bar ON bar.id = $bar_id
WHERE owner = $owner;

-- name: JoinNoConstraints :many
SELECT foo.*
FROM foo
CROSS JOIN bar
WHERE bar.id = $bar_id AND owner = $owner;
