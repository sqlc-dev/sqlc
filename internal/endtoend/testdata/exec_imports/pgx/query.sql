CREATE TABLE foo (bar int, bars int[] not null);

-- name: Bar :exec
SELECT bar
FROM foo;

-- name: Bars :exec
SELECT bars
FROM foo;
