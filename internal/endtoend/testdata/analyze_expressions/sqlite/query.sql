-- name: PostStats :one
SELECT count(*) AS total, max(id) AS latest, min(created) AS first
FROM posts WHERE user_id = ?;

-- name: UserScores :one
SELECT avg(score) AS mean, sum(score) AS sum, group_concat(name) AS names
FROM users;

-- name: ListUsers :many
SELECT id, lower(name) AS lname, id + 1 AS next
FROM users WHERE id IN (?, ?);

-- name: ListPosts :many
SELECT id, title FROM posts ORDER BY created LIMIT ? OFFSET ?;

-- name: UserPosts :many
SELECT u.name, p.title, p.created
FROM users u LEFT JOIN posts p ON p.user_id = u.id
WHERE u.name LIKE ? || '%';

-- name: CreatePost :one
INSERT INTO posts (user_id, title, created) VALUES (?, ?, datetime('now'))
RETURNING id, created;
