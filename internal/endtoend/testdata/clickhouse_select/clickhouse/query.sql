-- name: ListEvents :many
SELECT id, name, tag, amount, created FROM events;
