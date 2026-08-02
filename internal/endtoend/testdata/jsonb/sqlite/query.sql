-- name: InsertFoo :exec
INSERT INTO foo (
    a,
    b,
    c,
    d,
    e,
    f,
    g,
    h
) VALUES (
    @a,
    @b,
    @c,
    @d,
    @e,
    @f,
    @g,
    @h
) RETURNING *;

-- name: SelectFoo :exec
SELECT * FROM foo;
