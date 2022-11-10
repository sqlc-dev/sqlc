CREATE TABLE foo (
    a json not null,
    b jsonb not null,
    c json,
    d jsonb
);

-- name: SelectFoo :exec
SELECT * FROM foo;
