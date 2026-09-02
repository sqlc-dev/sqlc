-- name: Aggregates :one
SELECT count() AS n, max(amount) AS top, sum(amount), min(tag) AS first_tag, uniq(name)
FROM events;

-- name: Expressions :many
SELECT id + 1 AS next_id, nullIf(name, '') AS maybe, coalesce(tag, 'none') AS tag_or_none,
       lower(name), toDate(now()) AS today, id > 1 AS big, if(id = 1, 'one', NULL) AS word
FROM events;

-- name: LeftJoinDefaults :many
SELECT e.id, e.name, u.email
FROM events e
LEFT JOIN users u ON u.id = e.id;

-- name: Positional :many
SELECT id, name FROM events WHERE name = ? AND amount > ? AND tag = ?;

-- name: Named :many
SELECT id, name FROM events
WHERE name = sqlc.arg(name) AND tag = sqlc.narg(tag) AND (amount > sqlc.arg(amount) OR amount < sqlc.arg(amount));

-- name: Functions :many
SELECT id FROM events WHERE lower(name) = ? AND id IN (?) AND toDate(now()) > ?;

-- name: Paging :many
SELECT id FROM events ORDER BY id LIMIT ? OFFSET ?;

-- name: Projected :one
SELECT ? AS echo, 'literal' AS lit;
