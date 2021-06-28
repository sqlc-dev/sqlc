CREATE TABLE bar (id serial not null, name text);

-- name: GetBars :many
SELECT DISTINCT ON (a.id) a.*
FROM bar a;
