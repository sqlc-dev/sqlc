CREATE EXTENSION IF NOT EXISTS hstore;

CREATE TABLE foo (
        bar hstore NOT NULL,
        baz hstore
);

-- name: ListBar :many
SELECT bar FROM foo;

-- name: ListBaz :many
SELECT baz FROM foo;


