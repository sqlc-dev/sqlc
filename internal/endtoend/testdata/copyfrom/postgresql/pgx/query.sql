CREATE TABLE foo (a text, b integer);

-- name: InsertValues :copyFrom
INSERT INTO foo (a, b) VALUES ($1, $2);

-- name: InsertSingleValue :copyFrom
INSERT INTO foo (a) VALUES ($1);
