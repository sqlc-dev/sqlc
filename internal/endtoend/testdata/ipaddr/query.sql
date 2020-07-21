CREATE TABLE foo (bar bool not null, "inet" inet not null, "cidr" cidr not null);

-- name: Get :many
SELECT bar, "inet", "cidr" FROM foo LIMIT $1;
