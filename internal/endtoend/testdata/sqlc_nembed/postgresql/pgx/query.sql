-- name: LeftJoin :many
SELECT sqlc.embed(users), sqlc.nembed(posts) FROM users
LEFT JOIN posts ON users.id = posts.user_id;

-- name: LeftJoinOne :one
SELECT sqlc.embed(users), sqlc.nembed(posts) FROM users
LEFT JOIN posts ON users.id = posts.user_id
LIMIT 1;

-- name: NembedOnly :one
SELECT sqlc.nembed(users) FROM users WHERE id = $1;

-- name: WithAlias :many
SELECT sqlc.embed(u), sqlc.nembed(p) FROM users u
LEFT JOIN posts p ON u.id = p.user_id;
