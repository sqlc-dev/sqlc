CREATE TABLE bar (id serial not null);

-- name: In :many
SELECT *
FROM bar
WHERE id IN ($1, $2);
