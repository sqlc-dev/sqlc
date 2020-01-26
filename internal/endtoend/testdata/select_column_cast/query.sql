CREATE TABLE foo (bar bool not null);

-- name: SelectColumnCast :many
SELECT bar::int FROM foo;
