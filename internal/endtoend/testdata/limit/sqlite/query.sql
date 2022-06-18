CREATE TABLE foo (bar bool not null);

-- name: LimitMe :many
SELECT bar FROM foo LIMIT ?;
