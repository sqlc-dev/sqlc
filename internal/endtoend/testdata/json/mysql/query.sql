-- name: SelectFoo :exec
SELECT * FROM foo;

-- name: BulkInsert :copyfrom
INSERT INTO foo (a, b) VALUES (?, ?);
