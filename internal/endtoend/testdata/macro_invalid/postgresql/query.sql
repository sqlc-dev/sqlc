-- name: ByID :macro
-- @param id bigint
SELECT * FROM foo WHERE id = sqlc.arg(id);

-- name: UnknownMacro :many
SELECT * FROM sqlc.macro.Missing(1) m;

-- name: TooFewArguments :many
SELECT * FROM sqlc.macro.ByID() m;

-- name: TooManyArguments :many
SELECT * FROM sqlc.macro.ByID(1, 2) m;

-- name: NoArgumentList :many
SELECT * FROM sqlc.macro.ByID m;

-- name: Ping :macro
SELECT * FROM sqlc.macro.Pong() p;

-- name: Pong :macro
SELECT * FROM sqlc.macro.Ping() p;

-- name: Recursive :many
SELECT * FROM sqlc.macro.Ping() p;

-- name: ByID :macro
SELECT 1;

-- name: Nullable :macro
-- @param id
SELECT * FROM foo WHERE id = sqlc.narg(id);

-- name: NargOfParameter :many
SELECT * FROM sqlc.macro.Nullable(1) n;

-- name: Positional :macro
SELECT * FROM foo WHERE id = $1;

-- name: PositionalInMacro :many
SELECT * FROM sqlc.macro.Positional() p;

-- name: Twice :macro
-- @param id
-- @param ID
-- @param
SELECT sqlc.arg(id);

-- name: Empty :macro
;

-- A macro is only expanded where it is called, so this one never fails.
-- name: Unused :macro
SELECT * FROM sqlc.macro.Missing($1);

-- name: Valid :many
SELECT * FROM sqlc.macro.ByID($1) m;
