-- name: ListUsers :many
SELECT * FROM users;

-- name: CountUsers :one
SELECT count(*) AS total FROM users;

-- name: ListUserPosts :many
SELECT u.name, p.* FROM users u JOIN posts p ON p.user_id = u.id WHERE u.name = ?;

-- name: SearchPosts :many
SELECT id, title FROM posts WHERE title LIKE ? AND created > ?;

-- name: LatestPostPerUser :many
WITH latest AS (
  SELECT user_id, max(created) AS created FROM posts GROUP BY user_id
)
SELECT u.name, latest.created FROM users u JOIN latest ON latest.user_id = u.id;
