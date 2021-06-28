-- name: MissingSelect :many
SELECT * FROM nonexisting_relation WHERE name = 'foo';

-- name: MissingParam :many
SELECT * FROM nonexisting_relation WHERE name = $1;

-- name: MissingUpdate :exec
UPDATE nonexisting_relation SET name = $1;

-- name: MisingInsert :exec
INSERT INTO nonexisting_relation (id) VALUES ($1);
