-- name: ListUsers :many
SELECT * FROM users;

-- name: CountUsers :one
SELECT count(*) AS total FROM users;

-- name: ListUserPosts :many
SELECT u.name, p.* FROM users u JOIN posts p ON p.user_id = u.id WHERE u.name = $1;

-- name: TopUsers :many
WITH counts AS (
  SELECT user_id, count(*) AS n FROM posts GROUP BY user_id
)
SELECT u.name, c.n FROM users u JOIN counts c ON c.user_id = u.id
ORDER BY c.n DESC
LIMIT 5;
