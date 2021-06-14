CREATE SCHEMA foo;
CREATE TABLE foo.bar (id serial not null, name text not null);

-- name: SchemaScopedCreate :one
INSERT INTO foo.bar (id, name) VALUES ($1, $2) RETURNING id;
