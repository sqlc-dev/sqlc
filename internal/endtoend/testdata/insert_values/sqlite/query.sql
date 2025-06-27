/* name: InsertValues :exec */
INSERT INTO foo (a, b) VALUES (?, ?);

/* name: InsertMultipleValues :exec */
INSERT INTO foo (a, b) VALUES (?, ?), (?, ?);

/* name: InsertDefaultValues :exec */
INSERT INTO foo DEFAULT VALUES;
