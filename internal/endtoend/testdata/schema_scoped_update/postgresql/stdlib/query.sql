-- name: SchemaScopedUpdate :exec
UPDATE foo.bar SET name = $2 WHERE id = $1;
