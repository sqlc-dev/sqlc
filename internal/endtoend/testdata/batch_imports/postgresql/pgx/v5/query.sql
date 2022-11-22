CREATE SCHEMA myschema;
CREATE TABLE myschema.foo (a text, b integer);

-- name: InsertValues :batchone
INSERT INTO myschema.foo (a, b)
VALUES ($1, $2)
RETURNING a;

-- name: GetValues :batchmany
SELECT *
FROM myschema.foo
WHERE b = $1;

-- name: UpdateValues :exec
UPDATE myschema.foo SET a = $1, b = $2;
