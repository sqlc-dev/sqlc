-- name: Known :macro
SELECT 1;

-- name: Bad :one
SELECT * FROM sqlc.macro.Unknown(sqlc.arg(id)) u;

-- name: Good :one
SELECT * FROM sqlc.macro.Known() k WHERE id = sqlc.arg(id);

-- The name of a macro is matched exactly.
-- name: WrongCase :one
SELECT * FROM sqlc.macro.known() k;
