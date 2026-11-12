-- name: BugDemo :exec
SELECT id, ARRAY_AGG(whatever.*) FROM worker;
