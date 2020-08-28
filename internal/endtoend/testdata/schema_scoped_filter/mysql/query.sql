CREATE SCHEMA foo;
CREATE TABLE foo.bar (id serial not null);

-- name: SchemaScopedFilter :many
SELECT * FROM foo.bar WHERE id = ?;
