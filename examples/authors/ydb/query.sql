-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = $p0 LIMIT 1;

-- name: ListAuthors :many 
SELECT * FROM authors ORDER BY name;

-- name: CreateOrUpdateAuthor :exec 
UPSERT INTO authors (id, name, bio) VALUES ($p0, $p1, $p2);

-- name: DeleteAuthor :exec 
DELETE FROM authors WHERE id = $p0;

-- name: DropTable :exec
DROP TABLE IF EXISTS authors;