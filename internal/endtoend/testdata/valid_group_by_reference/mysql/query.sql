-- name: ListAuthors :many
SELECT   id, name as full_name, bio
FROM     authors
GROUP BY full_name;

-- name: ListAuthorsIdenticalAlias :many
SELECT   id, name as name, bio
FROM     authors
GROUP BY name;

-- name: ListMetrics :many
SELECT time_bucket('15 days', time) AS bucket, city_name, AVG(temp_c)
FROM weather_metrics
WHERE DATE_SUB(NOW(), INTERVAL 6 MONTH)
GROUP BY bucket, city_name
ORDER BY bucket DESC;
