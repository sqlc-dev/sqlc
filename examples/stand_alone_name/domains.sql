-- name: GetByNo :one
SELECT * FROM domains where tag = $1;