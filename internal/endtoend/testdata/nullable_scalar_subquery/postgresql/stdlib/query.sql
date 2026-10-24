-- name: ThroughView :many
select id, review_id, reviews from book_view;

-- name: Inline :many
select (select r.id from reviews r where r.book_id = b.id limit 1) review_id,
  (select count(*) from reviews r where r.book_id = b.id) reviews,
  (select 1) one
from books b;
