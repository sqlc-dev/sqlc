CREATE TABLE bar (tags text[] not null);

-- name: TextArray :many
SELECT * FROM bar;
