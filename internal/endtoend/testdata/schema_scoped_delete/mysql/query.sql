CREATE SCHEMA foo;
CREATE TABLE foo.bar (id serial not null);

-- name: SchemaScopedDelete :exec
DELETE FROM foo.bar WHERE id = ?;
