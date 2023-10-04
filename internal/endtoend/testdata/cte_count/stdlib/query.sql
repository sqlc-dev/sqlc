-- name: CTECount :many
WITH all_count AS (
	SELECT count(*) FROM bar
), ready_count AS (
	SELECT count(*) FROM bar WHERE ready = true
)
SELECT all_count.count, ready_count.count
FROM all_count, ready_count;
