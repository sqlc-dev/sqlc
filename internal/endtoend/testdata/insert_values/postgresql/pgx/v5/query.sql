-- name: InsertValues :exec
INSERT INTO foo (a, b) VALUES ($1, $2);

/* name: InsertMultipleValues :exec */
INSERT INTO foo (a, b) VALUES ($1, $2), ($3, $4);