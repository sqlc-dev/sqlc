-- name: SchemaScopedDelete :exec
DELETE FROM foo.bar WHERE id = $1;
