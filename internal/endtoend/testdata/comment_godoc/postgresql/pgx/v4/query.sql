CREATE TABLE foo (
  bar  text
);

-- name: ManyFoo :many
-- This function returns a list of Foos
SELECT * FROM foo;

-- name: OneFoo :one
-- This function returns one Foo
SELECT * FROM foo;

-- name: ExecFoo :exec
-- This function creates a Foo via :exec
INSERT INTO foo (bar) VALUES ("bar");

-- name: ExecRowFoo :execrows
-- This function creates a Foo via :execrows
INSERT INTO foo (bar) VALUES ("bar");

-- name: ExecResultFoo :execresult
-- This function creates a Foo via :execresult
INSERT INTO foo (bar) VALUES ("bar");
