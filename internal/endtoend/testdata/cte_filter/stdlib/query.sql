-- name: CTEFilter :many
WITH filter_count AS (
	SELECT count(*) FROM bar WHERE ready = $1
)
SELECT filter_count.count
FROM filter_count;
