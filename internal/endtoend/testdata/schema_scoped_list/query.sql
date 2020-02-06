CREATE SCHEMA foo;
CREATE TABLE foo.bar (id serial not null);

-- name: SchemaScopedList :many
SELECT * FROM foo.bar;
