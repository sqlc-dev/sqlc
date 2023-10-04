-- name: SchemaScopedFilter :many
SELECT * FROM foo.bar WHERE id = $1;
