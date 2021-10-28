CREATE TABLE bar (id bigserial not null);

-- name: Any :many
SELECT id
FROM bar
WHERE foo = ANY($1::bigserial[]);
