-- name: InnerWithDoesNotLeak :one
WITH x AS (SELECT 1 AS n)
SELECT
  (SELECT count(*) FROM (WITH authors AS (SELECT 1 AS id) SELECT id FROM authors) s) AS inner_count,
  (SELECT name FROM authors LIMIT 1) AS outer_name
FROM x;
