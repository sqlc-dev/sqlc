-- name: GetUser :one
WITH found AS (
    SELECT id, id AS id2 FROM users WHERE id = $1
)
SELECT id2 + $2 AS result FROM found;
