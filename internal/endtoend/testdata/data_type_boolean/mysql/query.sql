CREATE TABLE foo
(
    col_a BOOL             NOT NULL,
    col_b BOOLEAN          NOT NULL,
    col_c TINYINT(1)       NOT NULL
);

-- name: ListFoo :many
SELECT * FROM foo;

CREATE TABLE bar
(
    col_a BOOL,
    col_b BOOLEAN,
    col_c TINYINT(1)
);

-- name: ListBar :many
SELECT * FROM bar;
