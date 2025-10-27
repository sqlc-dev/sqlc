-- name: AuthorPages :many
SELECT author, Count(title) AS num_books, CAST(Sum(pages) AS Int32) AS total_pages
FROM books
GROUP BY author;


