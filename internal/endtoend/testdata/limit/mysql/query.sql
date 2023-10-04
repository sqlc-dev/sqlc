-- name: LimitMe :exec
UPDATE foo SET bar='baz' LIMIT ?;

-- name: LimitMeToo :exec
DELETE FROM foo LIMIT ?;
