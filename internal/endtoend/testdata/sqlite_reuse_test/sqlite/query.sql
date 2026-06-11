-- name: ReuseParam :exec
UPDATE foo SET name = @name WHERE id = @id OR name = @name;
