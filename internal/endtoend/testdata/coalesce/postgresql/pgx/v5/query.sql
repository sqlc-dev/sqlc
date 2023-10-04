-- name: CoalesceString :many
SELECT coalesce(bar, '') as login
FROM foo;

-- name: CoalesceNumeric :many
SELECT coalesce(baz, 0) as login
FROM foo;

-- name: CoalesceStringColumns :many
SELECT bar, bat, coalesce(bar, bat)
FROM foo;

-- name: CoalesceNumericColumns :many
SELECT baz, qux, coalesce(baz, qux)
FROM foo;

-- name: CoalesceStringNull :many
SELECT bar, coalesce(bar)
FROM foo;

-- name: CoalesceNumericNull :many
SELECT baz, coalesce(baz)
FROM foo;
