-- name: AdvisoryLock :many
SELECT pg_advisory_xact_lock($1);
