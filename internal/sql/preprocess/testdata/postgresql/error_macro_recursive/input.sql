-- name: Self :macro
SELECT * FROM sqlc.macro.Self() s;

-- name: Ping :macro
SELECT * FROM sqlc.macro.Pong() p;

-- name: Pong :macro
SELECT * FROM sqlc.macro.Ping() p;

-- name: Direct :many
SELECT * FROM sqlc.macro.Self() s;

-- name: Mutual :many
SELECT * FROM sqlc.macro.Ping() p;
