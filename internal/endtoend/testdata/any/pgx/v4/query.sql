-- name: Any :many
SELECT id
FROM bar
WHERE foo = ANY($1::bigserial[]);
