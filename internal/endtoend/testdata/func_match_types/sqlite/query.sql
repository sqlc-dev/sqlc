-- name: AuthorPages :many
SELECT author, count(title) AS num_books, sum(pages) AS total_pages
FROM books
GROUP BY author;
