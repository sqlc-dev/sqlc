CREATE TABLE bar (id INTEGER PRIMARY KEY, name text);

-- name: GetBars :many
SELECT DISTINCT * FROM bar;
