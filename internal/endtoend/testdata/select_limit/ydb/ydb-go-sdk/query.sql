-- name: FooLimit :many
SELECT a FROM foo
LIMIT $limit;

-- name: FooLimitOffset :many
SELECT a FROM foo
LIMIT $limit OFFSET $offset;

