CREATE SCHEMA myschema;
CREATE TABLE myschema.foo (a text, b integer);

-- name: InsertValues :copyfrom
INSERT INTO myschema.foo (a, b) VALUES ($1, $2);

-- name: InsertSingleValue :copyfrom
INSERT INTO myschema.foo (a) VALUES ($1);
