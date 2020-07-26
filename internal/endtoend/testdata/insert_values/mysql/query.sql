CREATE TABLE foo (a text, b integer);

/* name: InsertValues :exec */
INSERT INTO foo (a, b) VALUES (?, ?);
