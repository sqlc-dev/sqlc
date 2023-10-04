/* name: InsertSelect :exec */
INSERT INTO foo (name, meta)
SELECT name, ?
FROM bar WHERE ready = ?;
