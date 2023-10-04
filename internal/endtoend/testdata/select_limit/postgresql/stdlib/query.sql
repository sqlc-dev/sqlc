-- name: FooLimit :many
SELECT a FROM foo
LIMIT $1;

-- name: FooLimitOffset :many
SELECT a FROM foo
LIMIT $1 OFFSET $2;
