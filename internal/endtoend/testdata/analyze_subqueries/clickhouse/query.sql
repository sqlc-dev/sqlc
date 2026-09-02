-- name: Aliased :many
SELECT s.x, s.total, s.email
FROM (
    SELECT e.id AS x, sum(e.amount) AS total, any(u.email) AS email
    FROM events e JOIN users u ON u.id = e.id
    GROUP BY e.id
) s;

-- name: Union :many
SELECT id, name FROM events
UNION ALL
SELECT id, email FROM users;

-- name: ScalarSubquery :many
SELECT id, (SELECT count() FROM users) AS user_count
FROM events
WHERE id IN (SELECT id FROM users WHERE email = ?);

