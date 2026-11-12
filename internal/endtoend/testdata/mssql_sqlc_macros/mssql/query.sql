-- name: CreateAuthor :exec
INSERT INTO authors (id, name, bio) VALUES (sqlc.arg(id), sqlc.arg(name), sqlc.narg(bio));

-- name: UpdateBio :execrows
UPDATE authors SET bio = sqlc.narg(bio) WHERE id = @id;

-- name: ListAuthorsByIDs :many
SELECT id, name FROM authors WHERE id IN (sqlc.slice(ids)) ORDER BY id;

-- name: ListAuthorsByIDsAndName :many
SELECT id, name FROM authors WHERE name <> @name AND id IN (sqlc.slice(ids)) AND bio IS NULL ORDER BY id;

-- name: RowCount :one
SELECT CAST(@@ROWCOUNT AS INT) AS row_count;

-- name: AuthorBooks :many
SELECT sqlc.embed(a), sqlc.embed(books), books.title AS book_title
FROM authors a JOIN books ON books.author_id = a.id
WHERE a.id = sqlc.arg(author_id)
ORDER BY books.id;
