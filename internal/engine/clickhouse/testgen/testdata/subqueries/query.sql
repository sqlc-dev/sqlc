-- name: Cte :many
WITH t AS (SELECT id, tag FROM events)
SELECT t.id, t.tag, s.cnt
FROM t
JOIN (SELECT id, count() AS cnt FROM events GROUP BY id) s ON s.id = t.id;

-- name: Aliased :many
SELECT s.x, s.total, s.email
FROM (
    SELECT e.id AS x, sum(e.amount) AS total, any(u.email) AS email
    FROM events e JOIN users u ON u.id = e.id
    GROUP BY e.id
) s;

-- name: Star :many
SELECT * FROM users u JOIN events e ON e.id = u.id;

-- name: Union :many
SELECT id, name FROM events
UNION ALL
SELECT id, email FROM users;

-- name: ScalarSubquery :many
SELECT id, (SELECT count() FROM users) AS user_count
FROM events
WHERE id IN (SELECT id FROM users WHERE email = ?);
