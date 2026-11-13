-- name: Related :many
SELECT a.id, r.n
FROM authors a
CROSS APPLY (
    SELECT 1 AS n FROM books WHERE author_id = a.id
    UNION ALL
    SELECT 2 FROM authors WHERE id = a.id
) r;
