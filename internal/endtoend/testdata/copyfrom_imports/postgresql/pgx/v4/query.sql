-- name: InsertValues :copyfrom
INSERT INTO myschema.foo (a, b) VALUES ($1, $2);

-- name: InsertSingleValue :exec
INSERT INTO myschema.foo (a) VALUES ($1);

-- name: DeleteValues :execresult
DELETE
FROM myschema.foo;
