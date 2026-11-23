-- name: Clamp :macro
-- @param v
-- @param lo
-- @param hi
least(greatest(sqlc.arg(v), sqlc.arg(lo)), sqlc.arg(hi));

-- name: Clamped :many
SELECT sqlc.macro.Clamp(
  amount, -- the value
  /* the bounds */ sqlc.arg(lo),
  sqlc.arg(hi)
) FROM prices WHERE id = sqlc.arg(id);
