-- name: SchemaScopedCreate :one
INSERT INTO foo.bar (id, name) VALUES ($1, $2) RETURNING id;
