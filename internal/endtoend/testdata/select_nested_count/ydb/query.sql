-- name: GetAuthorsWithBooksCount :many
SELECT 
  a.id,
  a.name,
  a.bio,
  COUNT(b.id) AS books_count
FROM authors a
LEFT JOIN books b ON b.author_id = a.id
GROUP BY a.id, a.name, a.bio;
