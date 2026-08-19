-- name: CreateAuthor :one
INSERT INTO authors (name, bio) OUTPUT INSERTED.id VALUES (@name, @bio);

-- name: UpdateAuthorAlias :exec
UPDATE a SET name = @name FROM authors a WHERE a.id = @id;

-- name: UpdateBookPrices :exec
UPDATE b SET price = @price FROM books b JOIN authors a ON b.author_id = a.id WHERE a.name = @author;

-- name: DeleteBooksByAuthor :exec
DELETE b FROM books b INNER JOIN authors a ON b.author_id = a.id WHERE a.name = @name;
