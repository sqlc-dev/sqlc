-- name: DeleteOldLogs :exec
DELETE FROM logs WHERE created_at < ?;

-- name: DeleteErrorLogs :exec
DELETE FROM logs WHERE level = 'ERROR';

-- name: UpdateLogLevel :exec
ALTER TABLE logs UPDATE level = ? WHERE id = ?;
