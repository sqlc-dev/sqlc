ATTACH 'baz.db' AS baz;

CREATE TABLE users (
    id integer PRIMARY KEY,
    name text NOT NULL,
    age integer
);

CREATE TABLE posts (
    id integer PRIMARY KEY,
    user_id integer NOT NULL
);

CREATE TABLE baz.users (
    id integer PRIMARY KEY,
    name text NOT NULL
);


-- name: Only :one
SELECT sqlc.embed(users) FROM users;

-- name: OnlyCamel :one
SELECT sqlc.embed(Users) FROM Users;

-- name: WithAlias :one
SELECT sqlc.embed(u) FROM users AS u;

-- name: WithAliasCamel :one
SELECT sqlc.embed(U) FROM users AS U;

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
