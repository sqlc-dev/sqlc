-- name: Inner :macro
-- @param x
sqlc.arg(x) + 1;

-- name: Outer :macro
-- @param y
sqlc.macro.Inner(sqlc.arg(y)) * sqlc.macro.Inner(2);

-- A macro may be an argument to itself: an argument belongs to the caller.
-- name: Nested :one
SELECT sqlc.macro.Outer(sqlc.macro.Inner(sqlc.macro.Inner(sqlc.arg(n))));
