-- name: LimitMe :many
SELECT bar FROM foo LIMIT ?;

-- name: UpdateLimit :exec
UPDATE foo SET bar='baz' LIMIT ?;

-- name: DeleteLimit :exec
DELETE FROM foo LIMIT ?;
