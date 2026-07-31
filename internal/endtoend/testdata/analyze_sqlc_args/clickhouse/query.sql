-- name: GetEvent :one
SELECT id, name FROM events WHERE id = sqlc.arg(event_id);

-- name: FilterEvents :many
SELECT id, name FROM events WHERE name = sqlc.arg(name) AND amount > sqlc.narg(min_amount);
