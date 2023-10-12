-- name: InsertFoo :exec
INSERT INTO foo (bar)
SELECT 1, $1, $2;