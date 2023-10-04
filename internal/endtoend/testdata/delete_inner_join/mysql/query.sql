/* name: RemoveAllAuthorsFromTheGreatGatsby :exec */
DELETE author_book
FROM
  author_book
  INNER JOIN book ON book.id = author_book.book_id
WHERE
  book.title = 'The Great Gatsby';