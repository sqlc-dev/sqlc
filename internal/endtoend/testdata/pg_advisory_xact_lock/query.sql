-- name: AdvisoryLock :many
SELECT pg_advisory_unlock($1);
