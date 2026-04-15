-- name: InsertUser :exec
INSERT INTO Users (full_name, "EmailAddress", created_at)
VALUES (?, ?, ?);

-- name: SelectUsers :many
SELECT id, full_name, "EmailAddress", created_at
FROM Users;
