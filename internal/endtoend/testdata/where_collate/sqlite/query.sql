CREATE TABLE accounts (
    id TEXT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,

    UNIQUE (name COLLATE NOCASE)
);

-- name: GetAccountByName :one
SELECT * FROM accounts
WHERE name = ? COLLATE NOCASE
LIMIT 1;