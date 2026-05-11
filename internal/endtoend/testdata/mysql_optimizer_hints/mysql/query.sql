-- name: InlineHint :one
SELECT /*+ MAX_EXECUTION_TIME(1000) */ bar FROM foo LIMIT 1;

-- name: MultilineHint :one
SELECT
/*+ MAX_EXECUTION_TIME(1000) */
bar
FROM foo
LIMIT 1;
