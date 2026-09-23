-- name: FindEntry :one
SELECT id, value, code FROM entries
WHERE value = sqlc.arg(value) AND code = sqlc.arg(code);
