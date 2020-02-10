CREATE TABLE bar (id serial not null);

-- name: baz :one
SELECT foo FROM bar WHERE baz = $4;

-- name: bar
SELECT foo FROM bar WHERE baz = $1 AND baz = $3;

-- name: foo :one
SELECT foo FROM bar;

-- name: Named :many
SELECT id FROM bar WHERE id = $1 AND sqlc.arg(named) = true;

-- stderr
-- # package querytest
-- query.sql:4:1: could not determine data type of parameter $1
-- query.sql:7:1: could not determine data type of parameter $2
-- query.sql:10:8: column "foo" does not exist
-- query.sql:13:1: query mixes positional parameters ($1) and named parameters (sqlc.arg or @arg)
