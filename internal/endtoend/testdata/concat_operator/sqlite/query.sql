CREATE TABLE foo(bar TEXT, baz TEXT);

-- name: Concat :many
SELECT bar || ' ' || baz FROM foo;