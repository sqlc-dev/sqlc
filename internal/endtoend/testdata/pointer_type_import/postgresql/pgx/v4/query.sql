CREATE TABLE foo (bar date, baz uuid);

-- name: List :many
SELECT * FROM foo;

-- name: Find :one
SELECT bar FROM foo WHERE baz = $1;
