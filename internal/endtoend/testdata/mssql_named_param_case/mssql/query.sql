-- name: ByName :many
SELECT id FROM authors WHERE name = @Name OR bio = @name;
