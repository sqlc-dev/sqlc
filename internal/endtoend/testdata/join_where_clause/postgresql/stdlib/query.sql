-- name: JoinWhereClause :many
SELECT foo.*
FROM foo
JOIN bar ON bar.id = barid
WHERE owner = $1;

-- name: JoinParamWhereClause :many
SELECT foo.*
FROM foo
JOIN bar ON bar.id = $2
WHERE owner = $1;

-- name: JoinNoConstraints :many
SELECT foo.*
FROM foo
CROSS JOIN bar
WHERE bar.id = $2 AND owner = $1;