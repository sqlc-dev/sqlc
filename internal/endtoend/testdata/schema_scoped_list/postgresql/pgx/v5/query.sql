CREATE SCHEMA foo;
CREATE TABLE foo.bar (id serial not null);

-- name: SchemaScopedList :many
SELECT * FROM foo.bar;

-- name: SchemaScopedColList :many
SELECT foo.bar.id FROM foo.bar;