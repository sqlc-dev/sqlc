-- name: SchemaScopedList :many
SELECT * FROM foo.bar;

-- name: SchemaScopedColList :many
SELECT foo.bar.id FROM foo.bar;