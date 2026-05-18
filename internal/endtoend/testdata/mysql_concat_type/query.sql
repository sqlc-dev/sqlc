-- name: GetTickets :many
SELECT
    id,
    title,
    ticket_status,
    created_at
FROM ticket
WHERE (sqlc.narg(start_time) IS NULL OR created_at >= sqlc.narg(start_time))
  AND (sqlc.narg(title_prefix) IS NULL OR title LIKE CONCAT(sqlc.narg(title_prefix), '%'))
ORDER BY id ASC;
