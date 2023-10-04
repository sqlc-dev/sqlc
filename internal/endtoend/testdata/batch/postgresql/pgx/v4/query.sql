-- name: InsertValues :batchone
INSERT INTO myschema.foo (a, b)
VALUES ($1, $2)
RETURNING a;

-- name: GetValues :batchmany
SELECT *
FROM myschema.foo
WHERE b = $1;

-- name: UpdateValues :batchexec
UPDATE myschema.foo SET a = $1, b = $2;
