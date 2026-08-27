-- name: Coalesce :many
SELECT coalesce(bar, '') AS login
FROM foo;

-- name: CoalesceColumns :many
SELECT bar, bat, coalesce(bar, bat)
FROM foo;
