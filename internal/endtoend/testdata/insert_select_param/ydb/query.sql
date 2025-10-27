-- name: InsertSelect :exec
INSERT INTO authors (id, name, bio) 
SELECT $id, a.name, a.bio
FROM authors a
WHERE a.name = $name;
