-- name: Pair :macro
-- @param a
-- @param b
SELECT sqlc.arg(a), sqlc.arg(b);

-- name: TooFew :one
SELECT * FROM sqlc.macro.Pair(1) p;

-- name: TooMany :one
SELECT * FROM sqlc.macro.Pair(1, 2, 3) p;

-- name: Missing :one
SELECT * FROM sqlc.macro.Pair(1, ) p;

-- name: NoList :one
SELECT * FROM sqlc.macro.Pair p;

-- name: Unclosed :one
SELECT * FROM sqlc.macro.Pair(1, 2 p;
