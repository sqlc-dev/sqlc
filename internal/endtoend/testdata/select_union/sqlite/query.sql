-- name: SelectUnion :many
SELECT * FROM foo
UNION
SELECT * FROM foo;

-- name: SelectUnionWithLimit :many
SELECT * FROM foo
UNION
SELECT * FROM foo
LIMIT ? OFFSET ?;

-- name: SelectExcept :many
SELECT * FROM foo
EXCEPT
SELECT * FROM foo;

-- name: SelectIntersect :many
SELECT * FROM foo
INTERSECT
SELECT * FROM foo;

-- name: SelectUnionOther :many
SELECT * FROM foo
UNION
SELECT * FROM bar;