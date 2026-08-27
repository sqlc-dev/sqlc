-- name: DeleteAuthor :exec
DELETE FROM books AS b
WHERE NOT b.author IN (
    SELECT a.name
    FROM authors AS a
    WHERE a.age >= ?
  )
  AND NOT b.translator IN (
    SELECT t.name
    FROM translators AS t
    WHERE t.age >= ?
  )
  AND b.year <= ?;
