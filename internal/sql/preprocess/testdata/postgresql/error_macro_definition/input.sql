-- name: Dup :macro
SELECT 1;

-- name: Dup :macro
SELECT 2;

-- name: Params :macro
-- @param id
-- @param ID bigint
-- @param
SELECT sqlc.arg(id);

-- name: Empty :macro
-- nothing but a comment
;

-- The first definition is the one that is called.
-- name: Use :one
SELECT * FROM sqlc.macro.Dup() d;
