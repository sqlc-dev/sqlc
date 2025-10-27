-- name: InsertSelect :exec
INSERT INTO foo (name, meta)
SELECT name, $meta
FROM bar WHERE ready = $ready;
