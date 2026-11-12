-- name: Bracketed :many
SELECT [sqlc.arg(x)], [a;b], [@id], [x]]@y] FROM t WHERE id = @id;
