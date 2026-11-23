-- name: One :macro
SELECT 1;

-- sqlc.macro.One() in a comment is left alone.
-- name: Skipped :one
SELECT 'sqlc.macro.One()', "sqlc.macro.One()", $$ sqlc.macro.One() $$ /* sqlc.macro.One() */;
