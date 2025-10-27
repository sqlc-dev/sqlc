-- name: Only :one
SELECT sqlc.embed(users) FROM users;

-- name: WithAlias :one
SELECT sqlc.embed(u) FROM users u;

-- name: WithSubquery :many
SELECT sqlc.embed(users), Count(*) AS total_count FROM users;

-- name: WithAsterisk :one
SELECT sqlc.embed(users), * FROM users;

-- name: Duplicate :one
SELECT sqlc.embed(users), sqlc.embed(users) FROM users;

-- name: Join :one
SELECT sqlc.embed(users), sqlc.embed(posts) FROM posts
INNER JOIN users ON posts.user_id = users.id;
