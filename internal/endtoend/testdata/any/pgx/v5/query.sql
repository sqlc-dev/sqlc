-- name: Any :many
SELECT id
FROM bar
WHERE id = ANY($1::bigint[]);
