-- name: CreateUser :exec
INSERT INTO users (id, email) VALUES (?, ?);

-- name: CreateEvent :exec
INSERT INTO events VALUES (?, ?, ?, ?);

-- name: CreateEvents :exec
INSERT INTO events (id, name) VALUES (?, ?), (?, ?);

-- name: DropUsers :exec
TRUNCATE TABLE users;
