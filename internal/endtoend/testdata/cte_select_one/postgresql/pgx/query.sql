-- name: TestRecursive :one
WITH t1 AS (
    select 1 as foo
)
SELECT * FROM t1;