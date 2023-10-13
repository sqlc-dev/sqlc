-- name: Only :one
SELECT sqlc.embed(users) FROM users;

-- name: WithAlias :one
SELECT sqlc.embed(u) FROM users AS u;

-- name: WithSubquery :many
SELECT sqlc.embed(users), (SELECT count(*) FROM users) AS total_count FROM users;

-- name: WithAsterisk :one
SELECT sqlc.embed(users), * FROM users;

-- name: Duplicate :one
SELECT sqlc.embed(users), sqlc.embed(users) FROM users;

-- name: Join :one
SELECT sqlc.embed(u), sqlc.embed(p) FROM posts AS p
INNER JOIN users AS u ON p.user_id = u.users.id;

-- name: WithSchema :one
SELECT sqlc.embed(bu) FROM baz.users AS bu;

-- name: WithCrossSchema :many
SELECT sqlc.embed(u), sqlc.embed(bu) FROM users AS u
INNER JOIN baz.users bu ON u.id = bu.id;
