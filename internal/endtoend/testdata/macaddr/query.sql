CREATE TABLE foo (bar bool not null, addr macaddr not null);

-- name: Get :many
SELECT bar, addr FROM foo LIMIT $1;
