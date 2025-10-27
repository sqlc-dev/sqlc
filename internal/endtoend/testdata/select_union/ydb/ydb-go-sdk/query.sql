-- name: SelectUnion :many
SELECT * FROM foo
UNION
SELECT * FROM foo;

-- name: SelectUnionWithLimit :many
SELECT * FROM foo
UNION
SELECT * FROM foo
LIMIT $limit OFFSET $offset;

-- name: SelectUnionOther :many
SELECT * FROM foo
UNION
SELECT * FROM bar;

-- name: SelectUnionAliased :many
(SELECT * FROM foo)
UNION
SELECT * FROM foo;



