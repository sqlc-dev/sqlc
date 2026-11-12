-- name: CreateBook :one
INSERT INTO books (id, author_id, title) OUTPUT INSERTED.* VALUES (@id, @author_id, @title);

-- name: RetitleByAuthor :many
UPDATE b SET title = @title
OUTPUT INSERTED.id, DELETED.title AS old_title
FROM books b JOIN authors a ON a.id = b.author_id
WHERE a.name = @name;

-- name: DeleteByAuthor :many
DELETE b OUTPUT DELETED.id, DELETED.title
FROM books b JOIN authors a ON a.id = b.author_id
WHERE a.name = @name;

-- name: DeleteBook :one
DELETE FROM books OUTPUT DELETED.* WHERE id = @id;
