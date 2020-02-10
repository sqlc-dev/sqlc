CREATE TABLE foo (email text not null);

-- name: FirstQuery :many
SELECT * FROM foo;

-- name: SecondQuery :many
SELECT * FROM foo WHERE email = $1

-- stderr
-- # package querytest
-- query.sql:7:1: missing semicolon at end of file
