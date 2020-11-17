-- name: AdvisoryLockOne :one
SELECT pg_advisory_lock($1);

-- name: AdvisoryUnlock :many
SELECT pg_advisory_unlock($1);

-- name: AdvisoryLockExecResult :execresult
SELECT pg_advisory_lock($1);
