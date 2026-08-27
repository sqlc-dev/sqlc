-- name: CountRecentEvents :one
WITH activity AS (
  SELECT happened_at AS activity_at
  FROM events
)
SELECT COUNT(*) FILTER (WHERE activity_at >= $1)::bigint AS recent_events
FROM activity;
