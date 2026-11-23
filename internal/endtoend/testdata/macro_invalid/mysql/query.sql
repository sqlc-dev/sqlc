-- name: ByID :macro
-- @param id signed
SELECT * FROM foo WHERE id = sqlc.arg(id);

-- name: UnknownMacro :many
SELECT * FROM sqlc.macro.Missing(1) m;

-- name: TooManyArguments :many
SELECT * FROM sqlc.macro.ByID(1, 2) m;

-- name: Positional :macro
SELECT * FROM foo WHERE id = ?;

-- name: PositionalInMacro :many
SELECT * FROM sqlc.macro.Positional() p;

-- name: Valid :many
SELECT * FROM sqlc.macro.ByID(?) m;
