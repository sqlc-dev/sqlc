-- name: GetAccountByName :one
SELECT *
FROM accounts
WHERE name = ? COLLATE nocase
LIMIT 1;
