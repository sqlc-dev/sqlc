-- name: AuthorsJSON :many
SELECT id, name FROM authors WHERE id > @min_id FOR JSON PATH;

-- name: AuthorsXML :many
SELECT id, name FROM authors FOR XML RAW;

-- name: AuthorsWithBooks :many
SELECT a.id,
    (SELECT b.id, b.title FROM books b WHERE b.author_id = a.id FOR JSON PATH) AS books,
    (SELECT b.title + N',' FROM books b WHERE b.author_id = a.id FOR XML PATH('')) AS titles,
    (SELECT b.id, b.title FROM books b WHERE b.author_id = a.id FOR XML PATH, TYPE) AS books_xml
FROM authors a
WHERE a.name = @name;
