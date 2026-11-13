-- name: TopAuthors :many
SELECT TOP (@n) id, name FROM authors ORDER BY id;

-- name: PageAuthors :many
SELECT id, name FROM authors ORDER BY id OFFSET @skip ROWS FETCH NEXT @take ROWS ONLY;
