-- name: InsertSelected :exec
INSERT INTO t (id, name, at, n) SELECT @a, @b, @c, @d;

-- name: CopyRow :exec
INSERT INTO t (id, name, at, n) SELECT @a, name, at, n + @inc FROM t WHERE id = @src;
