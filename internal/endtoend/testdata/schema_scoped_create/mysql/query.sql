CREATE SCHEMA foo;
CREATE TABLE foo.bar (id serial not null, name text not null);

-- name: SchemaScopedCreate :execresult
INSERT INTO foo.bar (id, name) VALUES (?, ?);
