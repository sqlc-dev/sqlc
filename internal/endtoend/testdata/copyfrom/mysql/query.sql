-- name: InsertValues :copyfrom
INSERT INTO foo (a, b, c, d) VALUES (?, ?, ?, ?);

-- name: InsertSingleValue :copyfrom
INSERT INTO foo (a) VALUES (?);
