-- name: GetUser :one
SELECT * FROM users WHERE id = ?;

-- name: ListUsers :many
SELECT * FROM users ORDER BY name;

-- name: CreateUser :exec
INSERT INTO users (id, name, email) VALUES (?, ?, ?);

-- name: GetUserPosts :many
SELECT * FROM posts WHERE user_id = ? ORDER BY created_at DESC;

-- name: CreatePost :exec
INSERT INTO posts (id, user_id, title, body) VALUES (?, ?, ?, ?);


