-- A typed parameter casts whatever it is given.
-- name: Between :macro
-- @param v numeric(10, 2)
-- @param lo double precision
-- @param hi
sqlc.arg(v) >= sqlc.arg(lo) AND sqlc.arg(v) <= sqlc.arg(hi);

-- name: InRange :many
SELECT * FROM prices WHERE sqlc.macro.Between(amount, $1, $2 + 1);
