-- name: ListAuthors :many
SELECT   id, name as name, bio
FROM     authors
GROUP BY name;

-- name: ListAuthorsIdenticalAlias :many
SELECT   id, name as name, bio
FROM     authors
GROUP BY name;

-- name: ListMetrics :many
SELECT time_bucket('15 days', time) AS bucket, city_name, AVG(temp_c)
FROM weather_metrics
WHERE time > NOW() - (6 * INTERVAL '1 month')
GROUP BY bucket, city_name
ORDER BY bucket DESC;
