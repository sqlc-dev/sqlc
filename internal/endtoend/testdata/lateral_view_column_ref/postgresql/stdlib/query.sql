-- name: GetFooLateral :many
SELECT val, result FROM foo_lateral;

-- name: GetFooLateralDirect :many
SELECT f.id, f.val, sub.result
FROM foo f
CROSS JOIN LATERAL (
  SELECT (f.val || '-direct')::text AS result
) sub;
