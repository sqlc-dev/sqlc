-- name: InsertVector :exec
INSERT INTO items (embedding) VALUES ($1);

-- name: NearestNeighbor :many
SELECT *
FROM items
ORDER BY embedding <-> $1
LIMIT 5;
