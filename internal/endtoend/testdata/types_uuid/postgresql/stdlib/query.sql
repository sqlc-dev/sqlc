CREATE TABLE foo (
    description text,
    bar uuid,
    baz uuid not null
);

-- name: List :many
SELECT * FROM foo;

-- name: Find :one
SELECT bar FROM foo WHERE baz = $1;
