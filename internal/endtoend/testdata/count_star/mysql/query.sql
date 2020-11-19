CREATE TABLE bar (id serial not null);

-- name: CountStar :one
SELECT count(*) FROM bar;
