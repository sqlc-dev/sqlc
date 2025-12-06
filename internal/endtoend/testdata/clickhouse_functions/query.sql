-- Date/Time functions
-- name: GetLogsForDateRange :many
SELECT
    id,
    ip_address,
    timestamp,
    url
FROM web_logs
WHERE toDate(timestamp) >= ? AND toDate(timestamp) <= ?
ORDER BY timestamp DESC;

-- name: GetLogsGroupedByHour :many
SELECT
    toStartOfHour(timestamp) as hour,
    COUNT(*) as request_count,
    AVG(response_time_ms) as avg_response_time
FROM web_logs
WHERE timestamp >= ?
GROUP BY hour
ORDER BY hour DESC;

-- String functions
-- name: SearchLogsByUrl :many
SELECT
    id,
    ip_address,
    url,
    response_time_ms
FROM web_logs
WHERE url LIKE ?
ORDER BY timestamp DESC
LIMIT ?;

-- name: GetHttpStatusCodes :many
SELECT
    status_code,
    COUNT(*) as count,
    AVG(response_time_ms) as avg_response_time
FROM web_logs
WHERE timestamp >= ?
GROUP BY status_code
ORDER BY count DESC;

-- Conditional expressions
-- name: GetSlowRequests :many
SELECT
    id,
    url,
    response_time_ms,
    CASE
        WHEN response_time_ms > 1000 THEN 'very_slow'
        WHEN response_time_ms > 500 THEN 'slow'
        WHEN response_time_ms > 100 THEN 'medium'
        ELSE 'fast'
    END as performance
FROM web_logs
WHERE response_time_ms > ?
ORDER BY response_time_ms DESC;

-- Type casting
-- name: GetLogsSummary :many
SELECT
    ip_address,
    COUNT(*) as request_count,
    CAST(AVG(response_time_ms) AS UInt32) as avg_time,
    CAST(MIN(timestamp) AS Date) as first_request,
    CAST(MAX(timestamp) AS Date) as last_request
FROM web_logs
WHERE timestamp >= ?
GROUP BY ip_address
ORDER BY request_count DESC;

-- Math functions
-- name: GetResponseTimeStats :many
SELECT
    status_code,
    COUNT(*) as count,
    round(AVG(response_time_ms), 2) as avg_response_time,
    sqrt(varPop(response_time_ms)) as stddev
FROM web_logs
WHERE timestamp >= ?
GROUP BY status_code
ORDER BY avg_response_time DESC;

-- Map type operations
-- name: GetMetricsByTag :many
SELECT
    name,
    value,
    tags,
    created_at
FROM metrics
WHERE tags['environment'] = ?
ORDER BY created_at DESC;
