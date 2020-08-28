CREATE TABLE foo (
    a json not null,
    b json
);

-- name: SelectFoo :exec
SELECT * FROM foo;
