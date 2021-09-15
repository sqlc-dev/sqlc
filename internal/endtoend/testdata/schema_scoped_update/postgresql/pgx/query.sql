CREATE SCHEMA foo;
CREATE TABLE foo.bar (id serial not null, name text not null);

-- name: SchemaScopedUpdate :exec
UPDATE foo.bar SET name = $2 WHERE id = $1;
