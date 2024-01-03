-- name: InsertFoo :exec
INSERT INTO foo (bar)
SELECT 1, ?, ?;
