CREATE TABLE foo (bar text);

-- name: InsertFoo :exec
INSERT INTO foo (bar)
SELECT 1, ?, ?;
