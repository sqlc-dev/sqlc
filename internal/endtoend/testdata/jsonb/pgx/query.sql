-- name: InsertFoo :exec
INSERT INTO foo (
    a,
    b,
    c,
    d
) VALUES (
    @a,
    @b,
    @c,
    @d
) RETURNING *;

-- name: SelectFoo :exec
SELECT * FROM foo;
