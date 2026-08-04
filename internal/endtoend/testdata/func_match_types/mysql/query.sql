-- name: AuthorPages :many
select author, count(title) as num_books, SUM(pages) as total_pages, SUM(score) AS sum_score, SUM(price) AS sum_price, SUM(avg_word_length) as sum_avg_length
from books
group by author;
