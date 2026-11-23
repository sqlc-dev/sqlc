-- A line break inside a literal cannot be taken out, so the literal cannot be
-- pasted without moving every line after it.
-- name: Note :macro
'two
lines';

-- name: Echo :macro
-- @param x
sqlc.arg(x);

-- name: Body :one
SELECT sqlc.macro.Note();

-- name: Argument :one
SELECT sqlc.macro.Echo(1), sqlc.macro.Echo($$ two
lines $$);
