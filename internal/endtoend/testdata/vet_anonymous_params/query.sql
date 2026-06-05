-- name: ListByIDs :many
SELECT id FROM bar
WHERE id = ANY($1::bigint[]);

-- name: GetByID :one
SELECT id FROM bar
WHERE id = $1;
