CREATE TABLE bar (id serial not null, name text);

-- name: GetBars :many
SELECT FROM bar LIMIT 5;