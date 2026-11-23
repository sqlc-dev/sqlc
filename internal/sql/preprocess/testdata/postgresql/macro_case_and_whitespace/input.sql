-- name: One :macro
-- @param X
sqlc.arg(x) + SQLC.ARG(X) + @X;

-- The sqlc.macro prefix is matched like every other sqlc function.
-- name: Spaced :one
SELECT SQLC.Macro.One(1), sqlc . macro . One ( 2 );
