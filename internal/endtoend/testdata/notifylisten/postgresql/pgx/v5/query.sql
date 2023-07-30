-- name: NotifyTest :exec
NOTIFY test;

-- name: NotifyWithMessage :exec
NOTIFY test, 'msg';

-- name: ListenTest :exec
LISTEN test;
