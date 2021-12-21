CREATE SCHEMA myschema;
CREATE TABLE myschema.foo (a text, b integer);

-- name: InsertValues :copyFrom
INSERT INTO myschema.foo (a, b) VALUES ($1, $2);

-- name: InsertSingleValue :copyFrom
INSERT INTO myschema.foo (a) VALUES ($1);
