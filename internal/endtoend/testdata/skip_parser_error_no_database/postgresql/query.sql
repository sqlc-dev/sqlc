-- name: GetTest :one
SELECT id FROM test WHERE id = $1;
