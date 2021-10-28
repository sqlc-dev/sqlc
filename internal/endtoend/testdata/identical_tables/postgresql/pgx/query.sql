CREATE TABLE foo (id text not null);
CREATE TABLE bar (id text not null);

-- name: IdenticalTable :many
SELECT * FROM foo;
