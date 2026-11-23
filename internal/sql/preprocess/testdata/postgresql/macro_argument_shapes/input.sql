-- An untyped parameter takes a single term as it is written and anything
-- else in parentheses.
-- name: Twice :macro
-- @param x
sqlc.arg(x) * 2;

-- name: Shapes :many
SELECT
  sqlc.macro.Twice(a),
  sqlc.macro.Twice(t.a),
  sqlc.macro.Twice("t"."a"),
  sqlc.macro.Twice(42),
  sqlc.macro.Twice(1.5),
  sqlc.macro.Twice('7'),
  sqlc.macro.Twice($1),
  sqlc.macro.Twice(@b),
  sqlc.macro.Twice(sqlc.arg(c)),
  sqlc.macro.Twice(abs(a)),
  sqlc.macro.Twice((a + 1)),
  sqlc.macro.Twice(a + 1),
  sqlc.macro.Twice(a::int),
  sqlc.macro.Twice(greatest(a, 1) + least(a, 2)),
  sqlc.macro.Twice(- a)
FROM t;
