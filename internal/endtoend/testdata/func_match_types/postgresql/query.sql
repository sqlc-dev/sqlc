-- name: AuthorPages :many
select author, count(title) as num_books, SUM(pages) as total_pages
from books
group by author;
