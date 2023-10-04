-- name: DeleteAuthor :exec
DELETE FROM
  books AS b
WHERE
  b.author NOT IN (
    SELECT
      a.name
    FROM
      authors a
    WHERE
      a.age >= ?
  )
  AND b.translator NOT IN (
    SELECT
      t.name
    FROM
      translators t
    WHERE
      t.age >= ?
  )
  AND b.year <= ?;