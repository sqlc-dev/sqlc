-- name: IterValues :iter
SELECT *
FROM myschema.foo
WHERE b = $1;
