CREATE TABLE foo (name text not null);

-- name: As :many
SELECT name, name AS "other_name" FROM foo;
