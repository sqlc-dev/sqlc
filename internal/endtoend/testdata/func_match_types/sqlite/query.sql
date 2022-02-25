-- name: AuthorPages :many
select author, count(title) as num_books, sum(pages) as total_pages
from books
group by author;
