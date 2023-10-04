-- name: GetAccountByName :one
SELECT * FROM accounts
WHERE name = ? COLLATE NOCASE
LIMIT 1;