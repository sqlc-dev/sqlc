-- name: InsertSelect :exec
INSERT INTO foo (name, meta)
SELECT name, $1
FROM bar WHERE ready = $2;
