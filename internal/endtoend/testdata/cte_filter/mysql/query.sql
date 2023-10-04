-- name: CTEFilter :many
WITH filter_count AS (
	SELECT count(*) FROM bar WHERE ready = ?
)
SELECT filter_count.count
FROM filter_count;
