-- name: LatestBook :many
SELECT a.id, a.name, b.title
FROM authors a
CROSS APPLY (SELECT TOP 1 title FROM books WHERE author_id = a.id ORDER BY id DESC) b;

-- name: TopBooks :many
SELECT a.name, b.title
FROM authors a
OUTER APPLY (SELECT TOP (@n) title FROM books WHERE author_id = a.id ORDER BY id) b
WHERE a.name = @name;
