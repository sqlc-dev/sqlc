-- Example queries for sqlc
CREATE TABLE authors (
    id   INT
);

CREATE TABLE books (
    id   INT,
    author_id INT,
    price INT
);

-- name: ListAuthorsByCheapestBook :many
SELECT
    author_id, min(b.price) as min_price
From books b inner join authors a on a.id = b.author_id
GROUP BY b.author_id
ORDER BY min_price;