CREATE TABLE foo (a text, b integer, c DATETIME, d DATE);

-- name: InsertValues :copyfrom
INSERT INTO foo (a, b, c, d) VALUES (?, ?, ?, ?);

-- name: InsertSingleValue :copyfrom
INSERT INTO foo (a) VALUES (?);
