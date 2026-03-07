-- name: ReuseWithSlice :many
SELECT * FROM mytable
WHERE typ IN (sqlc.slice(types)) AND (name = @name OR id = @id OR name = @name);
