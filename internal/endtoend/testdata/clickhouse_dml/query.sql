-- Insert operations
-- name: InsertLog :exec
INSERT INTO logs (id, level, message, timestamp, source)
VALUES (?, ?, ?, ?, ?);

-- name: InsertMultipleLogs :exec
INSERT INTO logs (id, level, message, timestamp, source)
VALUES
    (?, ?, ?, ?, ?),
    (?, ?, ?, ?, ?),
    (?, ?, ?, ?, ?);

-- name: InsertNotification :exec
INSERT INTO notifications (id, user_id, message, read_status, created_at)
VALUES (?, ?, ?, ?, ?);

-- Select for use in tests
-- name: GetLogs :many
SELECT id, level, message, timestamp, source
FROM logs
WHERE timestamp >= ?
ORDER BY timestamp DESC
LIMIT ?;

-- name: GetLogsByLevel :many
SELECT id, level, message, timestamp, source
FROM logs
WHERE level = ?
ORDER BY timestamp DESC;

-- name: GetNotifications :many
SELECT id, user_id, message, read_status, created_at
FROM notifications
WHERE user_id = ?
ORDER BY created_at DESC;

-- name: CountUnreadNotifications :one
SELECT COUNT(*) as unread_count
FROM notifications
WHERE user_id = ? AND read_status = 0;

-- name: GetNotificationSummary :many
SELECT
    level,
    COUNT(*) as count
FROM logs
WHERE timestamp >= ?
GROUP BY level
ORDER BY count DESC;
