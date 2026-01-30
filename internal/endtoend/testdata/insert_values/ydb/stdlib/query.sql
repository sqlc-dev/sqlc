-- name: InsertValues :exec
INSERT INTO foo (a, b) VALUES ($a, $b);

/* name: InsertMultipleValues :exec */
INSERT INTO foo (a, b) VALUES ($a1, $b1), ($a2, $b2);
