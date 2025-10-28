-- name: GetUser :one
SELECT id, name, email, created_at
FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT id, name, email, created_at
FROM users
ORDER BY name;

-- name: CreateUser :one
INSERT INTO users (name, email)
VALUES ($1, $2)
RETURNING id, name, email, created_at;

-- name: GetUserPosts :many
SELECT p.id, p.title, p.content, p.published, p.created_at
FROM posts p
WHERE p.user_id = $1
ORDER BY p.created_at DESC;

-- name: CreatePost :one
INSERT INTO posts (user_id, title, content, published)
VALUES ($1, $2, $3, $4)
RETURNING id, user_id, title, content, published, created_at;
