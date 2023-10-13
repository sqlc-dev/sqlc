-- name: InsertBarBaz :exec
INSERT INTO foo (bar, baz) VALUES ($1);

-- name: InsertBar :exec
INSERT INTO foo (bar) VALUES ($1, $2);