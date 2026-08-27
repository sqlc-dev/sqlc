-- name: GetRepro :one
SELECT * FROM repro WHERE id = ? LIMIT 1;
