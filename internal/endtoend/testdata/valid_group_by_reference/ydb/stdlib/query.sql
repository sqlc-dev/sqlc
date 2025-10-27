-- name: ListAuthors :many
SELECT id, name AS name, bio
FROM authors;

-- name: ListAuthorsIdenticalAlias :many
SELECT id, name AS name, bio
FROM authors;

-- name: ListMetrics :many
SELECT DateTime::Format("%Y-%m-%d")(time) AS bucket, city_name, Avg(temp_c)
FROM weather_metrics
WHERE time > DateTime::MakeTimestamp(DateTime::Now()) - Interval("P6M")
GROUP BY bucket, city_name
ORDER BY bucket DESC;

