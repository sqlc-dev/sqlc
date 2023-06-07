CREATE SCHEMA IF NOT EXISTS baz;

CREATE TABLE users (
    id integer NOT NULL PRIMARY KEY,
    name varchar(255) NOT NULL,
    age integer NULL
);

CREATE TABLE posts (
    id integer NOT NULL PRIMARY KEY,
    user_id integer NOT NULL
);

CREATE TABLE baz.users (
    id integer NOT NULL PRIMARY KEY,
    name varchar(255) NOT NULL
);


-- name: Only :one
SELECT sqlc.embed(users) FROM users;

-- name: WithAlias :one
SELECT sqlc.embed(u) FROM users u;

-- name: WithSubquery :many
SELECT sqlc.embed(users), (SELECT count(*) FROM users) AS total_count FROM users;

-- name: WithAsterisk :one
SELECT sqlc.embed(users), * FROM users;

-- name: Duplicate :one
SELECT sqlc.embed(users), sqlc.embed(users) FROM users;

-- name: Join :one
SELECT sqlc.embed(users), sqlc.embed(posts) FROM posts
INNER JOIN users ON posts.user_id = users.id;

-- name: WithSchema :one
SELECT sqlc.embed(bu) FROM baz.users bu;

-- name: WithCrossSchema :many
SELECT sqlc.embed(users), sqlc.embed(bu) FROM users
INNER JOIN baz.users bu ON users.id = bu.id;