-- name: CreateAuthor :one
INSERT INTO authors (id, name, bio) VALUES ($1, $2, $3) RETURNING id;

-- name: UpdateAuthor :exec
UPDATE authors SET name = $name, bio = $bio WHERE id = $id;

-- name: DeleteBooksByAuthor :exec
DELETE FROM books USING authors a WHERE books.author_id = a.id AND a.name = $1;

-- name: TitlesFromBooks :many
INSERT INTO books (id, author_id, title)
SELECT id + 1000, author_id, upper(title) FROM books WHERE price > $min_price
RETURNING id, title;
