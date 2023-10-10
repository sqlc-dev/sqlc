-- name: InsertSelect :exec
INSERT INTO authors (id, name, bio) 
SELECT @id, name, bio
FROM authors
WHERE name = @name;