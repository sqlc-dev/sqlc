-- name: GetAccount :one
SELECT * FROM accounts WHERE id = $1;

-- name: ListAccounts :many
SELECT * FROM accounts ORDER BY id;

-- name: CreateAccount :one
INSERT INTO accounts (name, balance) VALUES ($1, $2) RETURNING *;

-- name: UpdateBalance :exec
UPDATE accounts SET balance = $1 WHERE id = $2;
