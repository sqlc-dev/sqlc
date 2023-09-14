CREATE TABLE foo (barid serial not null);
CREATE TABLE bar (id serial not null, owner text not null);

-- name: JoinWhereClause :many
SELECT foo.*
FROM foo
JOIN bar ON bar.id = barid
WHERE owner = ?;

-- name: JoinParamWhereClause :many
SELECT foo.*
FROM foo
JOIN bar ON bar.id = ?
WHERE owner = ?;

-- name: JoinNoConstraints :many
SELECT foo.*
FROM foo
CROSS JOIN bar
WHERE bar.id = ? AND owner = ?;
