-- name: CopyFromAuthor :many
WITH a AS (SELECT id FROM authors WHERE name = @name)
INSERT INTO books (id, author_id, title)
OUTPUT INSERTED.id, INSERTED.title
SELECT CAST(@book_id AS INT), a.id, CAST(@title AS NVARCHAR(200)) FROM a;
