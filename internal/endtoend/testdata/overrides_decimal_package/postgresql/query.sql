-- name: GetPrice :one
SELECT * FROM prices WHERE id = $1;
