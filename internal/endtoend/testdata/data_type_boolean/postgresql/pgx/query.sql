CREATE TABLE foo
(
    col_a BOOL             NOT NULL,
    col_b BOOLEAN          NOT NULL
);

-- name: ListFoo :many
SELECT * FROM foo;

CREATE TABLE bar
(
    col_a BOOL,
    col_b BOOLEAN
);

-- name: ListBar :many
SELECT * FROM bar;
