CREATE TABLE foo (bar bool not null);

-- name: Limit :many
SELECT bar FROM foo LIMIT $1;
