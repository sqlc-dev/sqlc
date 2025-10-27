-- name: CoalesceString :many
SELECT COALESCE(bar, '') AS login
FROM foo;

-- name: CoalesceNumeric :many
SELECT COALESCE(baz, 0) AS login
FROM foo;

-- name: CoalesceStringColumns :many
SELECT bar, bat, COALESCE(bar, bat)
FROM foo;

-- name: CoalesceNumericColumns :many
SELECT baz, qux, COALESCE(baz, qux)
FROM foo;

-- name: CoalesceStringNull :many
SELECT bar, COALESCE(bar)
FROM foo;

-- name: CoalesceNumericNull :many
SELECT baz, COALESCE(baz)
FROM foo;
