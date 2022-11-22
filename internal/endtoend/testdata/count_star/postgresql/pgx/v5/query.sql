CREATE TABLE bar (id serial not null);

-- name: CountStarLower :one
SELECT count(*) FROM bar;

-- name: CountStarUpper :one
SELECT COUNT(*) FROM bar;
