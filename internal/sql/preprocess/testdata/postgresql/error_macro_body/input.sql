-- An error in a body is reported where the macro is called, and only there:
-- a macro nothing calls is never looked at.
-- name: Positional :macro
SELECT * FROM foo WHERE id = $1;

-- name: Nullable :macro
-- @param id
SELECT * FROM foo WHERE id = sqlc.narg(id);

-- name: Sliced :macro
-- @param ids
SELECT * FROM foo WHERE id = ANY(sqlc.slice(ids));

-- name: Unused :macro
SELECT * FROM sqlc.macro.Nowhere($1);

-- name: UsePositional :many
SELECT * FROM sqlc.macro.Positional() p;

-- name: UseNullable :many
SELECT * FROM sqlc.macro.Nullable(1) n;

-- name: UseSliced :many
SELECT * FROM sqlc.macro.Sliced(sqlc.arg(ids)) s;

-- An error in an argument is reported at the argument.
-- name: BadArgument :many
SELECT * FROM sqlc.macro.Nullable(1 + sqlc.macro.Nowhere()) n;
