-- A parameter the body uses twice is an argument written twice, and a ? is
-- sent each time it is written.
-- name: Between :macro
-- @param v
-- @param lo
-- @param hi
sqlc.arg(lo) <= sqlc.arg(v) AND sqlc.arg(v) <= sqlc.arg(hi);

-- name: InRange :many
SELECT * FROM prices WHERE sqlc.macro.Between(sqlc.arg(amount), 1, `max`);
