-- name: AdvisoryLockExec :exec
SELECT pg_advisory_lock($1);

-- name: AdvisoryLockExecRows :execrows
SELECT pg_advisory_lock($1);

-- name: AdvisoryLockExecResult :execresult
SELECT pg_advisory_lock($1);


